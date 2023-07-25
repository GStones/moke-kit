package common_tests

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"github.com/gstones/platform/services/common/utils"
)

type searchTest struct {
	Index string
	Value interface{}
}

func StoreCommonTest(store document.DocumentStore) (err error) {
	if tmp, err := utils.NewTempDir("provider_test"); err != nil {
		return err
	} else {
		defer tmp.Cleanup()

		type testCase struct {
			key         string
			value       string
			fieldPath   string
			expectedErr error
		}

		testCases := []testCase{
			{
				key:         "/test/keys",
				value:       `{"a":"b"}`,
				fieldPath:   "/keys",
				expectedErr: nil,
			},
			{
				key:         "/keys/to/test",
				value:       `{"b":"a"}`,
				fieldPath:   "/to",
				expectedErr: nil,
			},
			{
				key:         "/keys/to/test",
				value:       `{"a":{"b":{"c":true}}}`,
				fieldPath:   "/test",
				expectedErr: nil,
			},
			{
				key:         "/keys/to/test/stuff",
				value:       `{"Test": "Name", "Data": "TestData"}`,
				fieldPath:   "/stuff",
				expectedErr: nil,
			},
			{
				key:         "",
				value:       `{"a":"b"}`,
				fieldPath:   "/stuff",
				expectedErr: document.ErrInvalidKeyFormat,
			},
			{
				key:         "/json/stuff",
				value:       `{"a":"b"}`,
				fieldPath:   "/notAField",
				expectedErr: errors2.ErrMalformedData,
			},
		}

		// perform a quick query against an empty store
		dst := searchTest{}
		if amt, err := store.Scan(
			"test",
			document.WithDestination(&dst),
			document.MatchKeyValue("Index", "value"),
		); err != nil {
			return err
		} else if amt != 0 {
			return errors.New("Got response of more than zero objects on a scan of an empty store.")
		}

		// populate a pair of values
		a := &searchTest{
			Index: "a",
			Value: "a",
		}
		if _, err := store.Set(
			document.NewKey("/local/test_a"),
			document.WithSource(&a),
			document.WithAnyVersion(),
		); err != nil {
			return err
		}
		b := &searchTest{
			Index: "b",
			Value: "b",
		}
		bKey := document.NewKey("/local/test_b")
		if _, err := store.Set(bKey, document.WithSource(&b), document.WithAnyVersion()); err != nil {
			return err
		}

		// get 'b'
		if v, err := store.Get(bKey, document.WithDestination(&dst)); err != nil {
			return err
		} else if v == document.NoVersion {
			return errors.New("Get returned no version.")
		}

		var expectedVal string
		var expectedAmt int

		// scan for 'a'
		expectedVal = "a"
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dst),
			document.MatchKeyValue("Index", expectedVal),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v",
				expectedVal)
		} else if amt != expectedAmt {
			return errors.Errorf("Unable to find the expected value %v in the document store.", expectedVal)
		}

		// scan for 'c'
		expectedVal = "c"
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dst),
			document.MatchKeyValue("Index", expectedVal),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v",
				expectedVal)
		} else if amt == expectedAmt {
			return errors.Errorf("Scan for nonexistent index name %v returned: %v", expectedVal, dst)
		}

		// scan for 'keys' = 'a'
		expectedVal = "c"
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dst),
			document.MatchKeyValue("keys", expectedVal),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v",
				expectedVal)
		} else if amt == expectedAmt {
			return errors.Errorf("Scan for nonexistent index name %v returned: %v", expectedVal, dst)
		}

		// populate two more values for the batch result test
		key1 := "/local/test_c1"
		c1 := &searchTest{
			Index: "c",
			Value: "1",
		}
		key2 := "/local/test_c2"
		c2 := &searchTest{
			Index: "c",
			Value: "2",
		}
		if _, err := store.Set(document.NewKey(key1), document.WithSource(&c1), document.WithAnyVersion()); err != nil {
			return err
		}
		if _, err := store.Set(document.NewKey(key2), document.WithSource(&c2), document.WithAnyVersion()); err != nil {
			return err
		}

		// scan for 'c' now that it exists
		expectedVal = "c"
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dst),
			document.MatchKeyValue("Index", expectedVal),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v "+
				"now that it exists", expectedVal)
		} else if amt != expectedAmt {
			return errors.Errorf("Unable to find the expected value %v in the document store %v",
				expectedVal, store.Name())
		}

		// scan for 'c' into a size 1 array
		expectedVal = "c"
		expectedAmt = 1
		dstArr := [1]searchTest{}
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArr),
			document.MatchKeyValue("Index", expectedVal),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v"+
				" in a size 1 array", expectedVal)
		} else if amt == expectedAmt {
			if dstArr[0].Index != expectedVal {
				return errors.Errorf("Got bad object back from scanning a size 1 array: "+
					"%v expected, %v returned", expectedVal, dstArr[0].Index)
			}
		} else {
			return errors.Errorf("Unable to find the expected value %v in the document store %v",
				expectedVal, store.Name())
		}

		expectedVal = "c"
		expectedAmt = 2
		dstArr2 := [2]searchTest{}
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArr2),
			document.MatchKeyValue("Index", expectedVal),
			document.WithLimit(2),
		); err != nil {
			return errors.Wrapf(err, "Unexpected error encountered while scanning for value %v"+
				" in a size 2 array", expectedVal)
		} else if amt == expectedAmt {
			if dstArr2[0].Index != expectedVal {
				return errors.Errorf("Got bad object back in response: %v expected, %v returned",
					expectedVal, dstArr2[0].Index)
			} else if dstArr2[1].Index != expectedVal {
				return errors.Errorf("Got bad object back in response: %v expected, %v returned",
					expectedVal, dstArr2[1].Index)
			}
		} else {
			return errors.Errorf("Unable to find the expected value %v in the document store %v.",
				expectedVal, store.Name())
		}

		// Test regex scan
		dstArr3 := make([]searchTest, 3)
		regex := "2"
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArr3),
			document.MatchRegex(regex),
		); err != nil {
			return err
		} else if amt < 1 {
			return errors.Wrapf(err, "Unable to find entry matching regex '%v'"+regex)
		}

		// Test float scan
		f1 := &searchTest{
			Index: "one",
			Value: 1,
		}
		f2 := &searchTest{
			Index: "two",
			Value: 2.0,
		}
		f3 := &searchTest{
			Index: "three",
			Value: 3.14,
		}

		if _, err := store.Set(
			document.NewKey("/local/test_one"),
			document.WithSource(&f1),
			document.WithAnyVersion(),
		); err != nil {
			return err
		}
		if _, err := store.Set(
			document.NewKey("/local/test_two"),
			document.WithSource(&f2),
			document.WithAnyVersion(),
		); err != nil {
			return err
		}
		if _, err := store.Set(
			document.NewKey("/local/test_three"),
			document.WithSource(&f3),
			document.WithAnyVersion(),
		); err != nil {
			return err
		}

		dstArrF := make([]searchTest, 2)
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArrF),
			document.MatchNumber("Value", document.ScanOpEquals, 2),
			document.WithNoLimit(),
		); err != nil {
			return err
		} else if amt != expectedAmt {
			return errors.Errorf("Unable to find entry matching Value = 2.0")
		}

		dstArrF = make([]searchTest, 5)
		expectedAmt = 2
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArrF),
			document.MatchNumber("Value", document.ScanOpLessThan, 3),
			document.WithNoLimit(),
		); err != nil {
			return err
		} else if amt != expectedAmt {
			return errors.New("Unable to find 2 entries matching Value < 3.0")
		}

		dstArrF = make([]searchTest, 5)
		expectedAmt = 1
		if amt, err := store.Scan(
			"/local/",
			document.WithDestination(&dstArrF),
			document.MatchNumber("Value", document.ScanOpGreaterThan, 3.1),
			document.WithNoLimit(),
		); err != nil {
			return err
		} else if amt != expectedAmt {
			return errors.New("Unable to find entry matching Value > 3.1")
		}

		// Ensure that keys iteration works before we add stuff
		if list, err := store.ListKeys("/empty"); err != nil {
			return errors.Wrap(err, "Error encountered iterating over keys of an empty store:")
		} else if len(list) != 0 {
			return errors.Errorf("Encountered non-empty list of length %v when an empty list was expected.",
				len(list))
		}

		for i, tc := range testCases {
			var testKey document.Key
			if testKey, err = document.NewKeyFromString(tc.key); err != nil {
				if tc.expectedErr == err {
					continue
				} else {
					return errors.Wrapf(err, "Could not create new keys in test case #%v", i+1)
				}
			}

			watcher := TestDocWatcher{Logger: zap.NewNop()}
			store.AddWatcher(testKey, &watcher)

			var src interface{}
			if err = json.Unmarshal([]byte(tc.value), &src); err != nil {
				return errors.Wrapf(err, "Error encountered unmarshalling value to be set in test case #%v",
					i+1)
			}

			var version document.Version
			if version, err = store.Set(
				testKey,
				document.WithAnyVersion(),
				document.WithTTL(1*time.Minute),
				document.WithSource(src),
			); err != nil {
				return errors.Wrapf(err, "Error encountered setting the Document's value for the provided"+
					" keys %v in test case #%v", testKey.String(), i+1)
			} else if version == document.NoVersion {
				return errors.Wrapf(err, "Set() call did not set the version of the document in test case #%v",
					i+1)
			} else if ok, err := store.Contains(testKey); err != nil {
				return errors.Wrapf(err, "Error encountered in store.Contains() call in test case #%v",
					i+1)
			} else if !ok {
				return errors.Wrapf(err, "Document store %v does not contain the expected keys %v in test case"+
					"#%v", store.Name(), testKey.String(), i+1)
			}

			if version, err = store.SetField(
				testKey,
				tc.fieldPath,
				document.WithSource(src),
				document.WithVersion(version),
			); err != nil {
				return errors.Wrapf(err, "Could not set the field for the provided document in test case #%v",
					i+1)
			} else if version == document.NoVersion {
				return errors.Wrapf(err, "SetField() call did not set the version of the document "+
					"in test case #%v", i+1)
			}

			var dst interface{}
			if version, err = store.Get(
				testKey,
				document.WithDestination(&dst),
				document.WithVersion(version),
			); err != nil {
				err = errors.Wrapf(err, "Could not get the provided document in test case #%v", i+1)
				return err
			} else if version == document.NoVersion {
				return errors.Errorf("Get() call returned NoVersion in test case #%v", i+1)
			}

			type dupBlockCase struct {
				data string
			}

			dupBlockCases := []dupBlockCase{
				{"keys " + tc.key + "\nset test\n5\n*"},
				{"keys " + tc.key + "\nset test2\n0\n*"},
				{"keys " + tc.key + "\ncpy test test2\n"},
				{"keys " + tc.key + "\nset test2\n10\n*"},
				{"keys " + tc.key + "\nswp test2 test\n"},
				{"keys " + tc.key + "\nmov test2 test\n"},
				{"keys " + tc.key + "\ndel test2"},
			}

			var appliedDUPBlocks int
			for caseNum, dbc := range dupBlockCases {
				if reader, err := dupblock.NewTextReader(dupblock.WithBytes([]byte(dbc.data))); err != nil {
					return err
				} else if err := store.ApplyDUPBlock(reader); err != nil {
					return errors.Wrapf(err, "Test Case #%v: "+
						"DUPBlock was not successfully applied for DUPBlock case #%v", i+1, caseNum+1)
				}
				changed := watcher.DocChangeCalled() || watcher.DocDeleteCalled() || watcher.DocExpireCalled()
				if changed {
					appliedDUPBlocks++
					watcher.resetWatcherFlags(watcher.docUpdateMutex)
				} else {
					return errors.Wrapf(err, "Test Case #%v: "+
						"DUPBlock was not successfully applied for DUPBlock case #%v", i+1, caseNum+1)
				}
			}

			// Get the keys again, to make sure we have the current version after DUPBlocks are done
			if version, err = store.Get(
				testKey,
				document.WithDestination(&dst),
				document.WithVersion(version),
			); err != nil {
				return errors.Wrapf(err, "Could not get the provided document in test case #%v", i+1)
			} else if version == document.NoVersion {
				return errors.Wrapf(err, "Get() call returned NoVersion in test case #%v", i+1)
			}

			// Get the list of keys to ensure that we get the index expected
			prefix := document.KeySeparator + testKey.Parts()[0]
			if list, err := store.ListKeys(prefix); err != nil {
				return errors.Wrapf(err, "Error encountered listing keys in test case #%v", i+1)
			} else {
				for _, k := range list {
					if !strings.HasPrefix(k.String(), prefix) {
						return errors.Errorf("Unexpected keys %v returned in test case #%v, prefix of %v",
							k, i+1, prefix)
					}
				}
			}

			// Remove the keys, since we are done performing DUPBlock applications
			if err = store.Remove(
				testKey,
				document.WithSource(src),
				document.WithVersion(version),
			); err != nil {
				return errors.Wrapf(err, "Error encountered removing the keys %v with version %v "+
					"in test case #%v", testKey, version, i+1)
			}

			// Remove the keys again, to check that you cannot remove a nonexistent keys
			if err = store.Remove(
				testKey,
				document.WithSource(src),
				document.WithAnyVersion(),
			); err != nil && errors.Cause(err) != errors2.ErrKeyNotFound {
				return errors.Wrapf(err, "Error encountered removing the keys %v with any version "+
					"in test case #%v", testKey, i+1)
			}

			// Final look at keys listing
			if key, err := document.NewKeyFromString(tc.key); err != nil {
				if tc.expectedErr == err {
					continue
				} else {
					return errors.Wrapf(err, "Couldn't create new keys in post-test case #%v", i+1)
				}
			} else if _, err := store.Set(
				key,
				document.WithAnyVersion(),
				document.WithTTL(1*time.Minute),
				document.WithSource(tc),
			); err != nil {
				return errors.Wrapf(err, "Couldn't re-insert keys %v in post-test case #%v", key, i+1)
			}
			store.RemoveWatcher(testKey, &watcher)
		}
	}
	return nil
}

// Test Doc Watcher helper structs and funcs
type TestDocWatcher struct {
	key    document.Key
	store  document.DocumentStore
	Logger *zap.Logger

	docUpdateMutex  sync.Mutex
	docChangeCalled bool
	docDelCalled    bool
	docExpireCalled bool
}

func (w *TestDocWatcher) OnDocumentChanged(key document.Key) {
	w.docUpdateMutex.Lock()
	w.docChangeCalled = true
	w.docUpdateMutex.Unlock()
}

func (w *TestDocWatcher) OnDocumentDeleted(key document.Key) {
	w.docUpdateMutex.Lock()
	w.docDelCalled = true
	w.docUpdateMutex.Unlock()
}

func (w *TestDocWatcher) OnDocumentExpired(key document.Key) {
	w.docUpdateMutex.Lock()
	w.docExpireCalled = true
	w.docUpdateMutex.Unlock()
}

func (w *TestDocWatcher) DocChangeCalled() bool {
	w.docUpdateMutex.Lock()
	result := w.docChangeCalled
	w.docUpdateMutex.Unlock()

	return result
}

func (w *TestDocWatcher) DocDeleteCalled() bool {
	w.docUpdateMutex.Lock()
	result := w.docDelCalled
	w.docUpdateMutex.Unlock()

	return result
}

func (w *TestDocWatcher) DocExpireCalled() bool {
	w.docUpdateMutex.Lock()
	result := w.docExpireCalled
	w.docUpdateMutex.Unlock()

	return result
}

func (w *TestDocWatcher) resetWatcherFlags(mutex sync.Mutex) {
	mutex.Lock()
	w.docChangeCalled = false
	w.docExpireCalled = false
	w.docDelCalled = false
	mutex.Unlock()
}
