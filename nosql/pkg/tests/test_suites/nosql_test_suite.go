package test_suites

import (
	"encoding/json"
	"fmt"
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"runtime"
	"time"

	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/nosql/document/badger"
	"github.com/gstones/platform/services/common/nosql/document/couchbase"
	"github.com/gstones/platform/services/common/nosql/document/mock"
	"github.com/gstones/platform/services/common/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func TestSuite(impl string, config couchbase.ClusterConfig, documentStore string) {

	fmt.Println("		       ┌───────────────────────────────┐")
	fmt.Println("		       │        NoSQL Test Suite       │")
	fmt.Println("		       └───────────────────────────────┘")

	if impl == "" {
		_ = errEncountered(errors.New("Please specify a NoSQL implementation through the --impl flag " +
			"before running this test suite."))
		return
	}

	var err error

	// create the document store provider in question
	var provider document.DocumentStoreProvider
	switch impl {
	case "couchbase":
		fmt.Println("Creating a new Document Store Provider ...")
		if provider, err = couchbase.NewDocumentStoreProvider(config, zap.NewNop()); err != nil {
			_ = errEncountered(err)
			return
		}

	case "badger":
		if runtime.GOOS == "linux" {
			fmt.Println("NOTICE: If you are running this suite in a WSL environment,",
				"Badger is incompatible with WSL.")
		}

		fmt.Println("Badger test suite running with configuration:")
		fmt.Println("Document Store Name:", documentStore)
		var dir utils.TempDir
		if dir, err = utils.NewTempDir("nosql_test_suite_"); err != nil {
			_ = errEncountered(err)
			return
		}
		defer dir.Cleanup()

		fmt.Println("\nCreating a new Document Store Provider ...")
		if provider, err = badger.NewDocumentStoreProvider(dir.Path(), 5*time.Minute, zap.NewNop()); err != nil {
			_ = errEncountered(err)
			return
		}

	case "mock":
		fmt.Println("\nCreating a new Document Store Provider ...")
		if provider, err = mock.NewDocumentStoreProvider(); err != nil {
			_ = errEncountered(err)
			return
		}

	default:
		_ = errEncountered(errors.New("Invalid implementation specified. " +
			"Valid options are: couchbase, badger, mock."))
		return
	}

	fmt.Println("Successfully created a new Document Store Provider!")

	// open the provider and run the test suite against it
	fmt.Println("Opening the", documentStore, "Document Store ...")
	var store document.DocumentStore
	if store, err = provider.OpenDocumentStore(documentStore); err != nil {
		_ = errEncountered(err)
		return
	} else {
		fmt.Println("Successfully opened the", store.Name(), "Document Store!")

		// schedule a shutdown in case of errors
		defer func() {
			if provider == nil {
				return
			}
			if err = provider.Shutdown(); err != nil {
				_ = errEncountered(err)
			}
		}()

		if err = runCommonTests(store); err != nil {
			return
		}
	}

	// clean up before displaying a "pass"
	if err = provider.Shutdown(); err != nil {
		_ = errEncountered(err)
		return
	}
	// make sure the deferred cleanup doesn't panic
	provider = nil

	fmt.Println("		       ┌───────────────────────────────┐")
	fmt.Println("		       │    NoSQL Test Result: PASS    │")
	fmt.Println("		       └───────────────────────────────┘")
}

func runCommonTests(store document.DocumentStore) error {
	type testCase struct {
		data string
	}
	var err error

	fmt.Println("Creating a new Key with value '/test/keys' ...")
	var testKey document.Key
	if testKey, err = document.NewKeyFromString("/test/keys"); err != nil {
		return errEncountered(err)
	}
	fmt.Println("Successfully created a new Key with value '/test/keys'!")

	var testVersion document.Version
	var src interface{}
	if err = json.Unmarshal([]byte(`{"a":"b"}`), &src); err != nil {
		return errEncountered(err)
	}
	fmt.Println("Setting the value of", testKey.String(), "to", `{"a":"b"}`, "...")

	if testVersion, err = store.Set(
		testKey,
		document.WithAnyVersion(),
		document.WithTTL(1*time.Minute),
		document.WithSource(src)); err != nil {
		return errEncountered(err)
	} else {
		fmt.Println("Successfully set the value of", testKey.String(), "to", `{"a":"b"}`, "!")
		fmt.Println("Checking if the Document Store now contains the Key ...")
		if ok, err := store.Contains(testKey); err != nil {
			return errEncountered(err)
		} else if !ok {
			return errEncountered(errors.New("ErrKeyNotInStore"))
		}
		fmt.Println("The Document Store contains:", testKey.String())

		fmt.Println("Getting the value of", testKey.String(), "to confirm the keys's value was set ...")
		var dst interface{}
		if testVersion, err = store.Get(
			testKey,
			document.WithVersion(testVersion),
			document.WithDestination(&dst)); err != nil {
			return errEncountered(err)
		}
		fmt.Println("Value of", testKey.String(), "is:", dst)
		fmt.Println("Successfully set the value of keys", testKey.String()+"!")

		fmt.Println("Applying DUPBlock test cases to the document store ...")
		testCases := []testCase{
			{"keys /test/keys\nset test\n5\n*"},
			{"keys /test/keys\nset test2\n0\n*"},
			{"keys /test/keys\ncpy test test2\n"},
			{"keys /test/keys\nset test2\n10\n*"},
			{"keys /test/keys\nmov test2 test\n"},
		}
		for i, tc := range testCases {
			if reader, err := dupblock.NewTextReader(dupblock.WithBytes([]byte(tc.data))); err != nil {
				return errEncountered(err)
			} else if err := store.ApplyDUPBlock(reader); err != nil {
				return errEncountered(err)
			}
			fmt.Println("Test Case", i+1, "completed!")
		}
		fmt.Println("All DUPBlock test cases were applied!")

		// prep data for ListKeys...
		fmt.Println("Inserting dataset for ListKeys & Scan tests ...")
		for k, v := range map[string]string{
			"/test/list_test_1":      `{"Idx":"a","Index":"1"}`,
			"/test/list_test_2":      `{"Idx":"b","value":"2"}`,
			"/testkeys/list_test_3":  `{"Idx":"c","int":3}`,
			"/testkeys/list_test_4":  `{"Idx":"d","flag":false}`,
			"/testkeys/scan/test_5":  `{"Idx":"e","User":"Guybrush Threepwood"}`,
			"/testkeys/scan/test_6":  `{"Idx":"f","User":"Murray"}`,
			"/testkeys/scan/test_7":  `{"Idx":"g","User":"Murray"}`,
			"/testkeys/scan/test_8":  `{"Idx":"h","Value":1.0}`,
			"/testkeys/scan/test_9":  `{"Idx":"i","Value":2}`,
			"/testkeys/scan/test_10": `{"Idx":"j","Value":3.14}`,
		} {
			if err = json.Unmarshal([]byte(v), &src); err != nil {
				return errEncountered(err)
			}
			fmt.Printf("Setting the value of `%s` to `%s` ...\n", k, v)

			if _, err = store.Set(
				document.NewKeyFromStringUnchecked(k),
				document.WithAnyVersion(),
				document.WithTTL(1*time.Minute),
				document.WithSource(src),
			); err != nil {
				return errEncountered(err)
			}
		}

		// and quickly verify that the last keys from the set was added successfully
		if ok, err := store.Contains(document.NewKeyFromStringUnchecked("/testkeys/list_test_4")); err != nil {
			return errEncountered(err)
		} else if !ok {
			return errEncountered(errors.New("ErrKeyNotInStore"))
		}

		// test ListKeys
		fmt.Println("Listing keys that start with `/testkeys/` ...")
		if keyList, err := store.ListKeys("/testkeys/", document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if len(keyList) < 5 {
			err = errEncountered(errors.New("ErrUnexpectedResultLength"))
			fmt.Println(keyList)
			return err
		} else {
			fmt.Println("Successfully listed keys matching prefix `/testkeys/`:", keyList)
		}

		if keyList, err := store.ListKeys("/look-out-behind-you-it-s-a-three-headed-monkey/", document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if len(keyList) != 0 {
			err = errEncountered(errors.New("ErrUnexpectedResultLength"))
			fmt.Println(keyList)
			return err
		} else {
			fmt.Println("Successfully listed 0 keys matching nonexistent prefix.")
		}

		// validate ListKeys with limit and offset
		fmt.Println("Listing keys that start with `/testkeys/`, limit 3 ...")
		if keyList, err := store.ListKeys("/testkeys/", document.WithLimit(3)); err != nil {
			return errEncountered(err)
		} else if len(keyList) != 3 {
			err = errEncountered(errors.New("ErrUnexpectedResultLength"))
			fmt.Println(keyList)
			return err
		} else {
			fmt.Println("Successfully listed keys matching prefix `/testkeys/`:", keyList)
		}

		fmt.Println("Listing keys that start with `/testkeys/`, offset 3 limit 2 ...")
		if keyList, err := store.ListKeys("/testkeys/", document.WithOffset(3), document.WithLimit(2)); err != nil {
			return errEncountered(err)
		} else if len(keyList) != 2 {
			err = errEncountered(errors.New("ErrUnexpectedResultLength"))
			fmt.Println(keyList)
			return err
		} else {
			fmt.Println("Successfully listed keys matching prefix `/testkeys/`:", keyList)
		}

		// test Scan
		type scanTest struct {
			Idx   string
			User  string
			Value interface{}
		}
		scanDest := scanTest{}
		fmt.Println("Scanning for single existant user ...")
		reqUser := "Guybrush Threepwood"
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanDest), document.MatchKeyValue("User", reqUser)); err != nil {
			return errEncountered(err)
		} else if num != 1 {
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else if scanDest.User != reqUser {
			return errEncountered(errors.New("ErrResultMismatch"))
		} else {
			fmt.Printf("Found record for user %s.\n", reqUser)
		}

		fmt.Println("Scanning for nonexistant user ...")
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanDest), document.MatchKeyValue("User", "LeChuck")); err != nil {
			return errEncountered(err)
		} else if num != 0 {
			// we're not expecting a response, freak out
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		}

		fmt.Println("Scanning for multi-entry user ...")
		reqUser = "Murray"
		expectedRows := 2
		scanRes := make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchKeyValue("User", reqUser), document.WithLimit(expectedRows)); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for %s, got %d results:", expectedRows, reqUser, num)
			for _, row := range scanRes {
				fmt.Printf(" %s", row.Idx)
			}
			fmt.Println()
		}

		// test WithOffset
		fmt.Println("Scanning with offset...")
		expectedRows = 1
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchKeyValue("User", reqUser), document.WithOffset(1), document.WithLimit(expectedRows)); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there is one entry for this user, after offset
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for %s, got %d results:", expectedRows, reqUser, num)
			for _, row := range scanRes {
				fmt.Printf(" %s", row.Idx)
			}
			fmt.Println()
		}

		// test WithKeyLike
		fmt.Println("Scanning for users matching pattern...")
		expectedRows = 2
		reqUser = "Mu%"
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchKeyLike("User", reqUser), document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for %s, got %d results:", expectedRows, reqUser, num)
			for _, row := range scanRes {
				fmt.Printf(" %s", row.Idx)
			}
			fmt.Println()
		}

		// test unset scan type
		fmt.Println("Scanning with unset type...")
		expectedRows = 0
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes)); err == nil {
			return errEncountered(errors.Wrap(errors2.ErrInternal, "got nil instead of expected error"))
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else if err != errors2.ErrNoScanType && errors.Cause(err) != errors2.ErrNoScanType {
			return errEncountered(errors.Wrap(errors2.ErrInternal, err.Error()))
		}

		// test multiple chained queries
		fmt.Println("Scanning for multiple conditions ...")
		reqUser = "Murray"
		expectedRows = 1
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchKeyValue("User", reqUser), document.MatchKeyValue("Idx", "g"), document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for User = %s, Idx = g, got %d results:", expectedRows, reqUser, num)
			for _, row := range scanRes {
				if row.Idx != "" {
					fmt.Printf(" %s", row)
				}
			}
			fmt.Println()
		}

		// test multiple chained queries
		fmt.Println("Scanning for numerical equality ...")
		reqNum := 2.0
		expectedRows = 1
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchNumber("Value", document.ScanOpEquals, reqNum), document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for Value = %.0f, got %d results:", expectedRows, reqNum, num)
			for _, row := range scanRes {
				if row.Idx != "" {
					fmt.Printf(" %.2f", row.Value)
				}
			}
			fmt.Println()
		}

		fmt.Println("Scanning for numerical < ...")
		reqNum = 3
		expectedRows = 2
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchNumber("Value", document.ScanOpLessThan, reqNum), document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for Value < %.1f, got %d results:", expectedRows, reqNum, num)
			for _, row := range scanRes {
				if row.Idx != "" {
					fmt.Printf(" %.2f", row.Value)
				}
			}
			fmt.Println()
		}

		fmt.Println("Scanning for numerical > ...")
		reqNum = 3.1
		expectedRows = 1
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchNumber("Value", document.ScanOpGreaterThan, reqNum), document.WithNoLimit()); err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for Value > %.1f, got %d results:", expectedRows, reqNum, num)
			for _, row := range scanRes {
				if row.Idx != "" {
					fmt.Printf(" %.2f", row.Value)
				}
			}
			fmt.Println()
		}

		// test MatchRegex
		regex := "e+pw.o"
		fmt.Printf("Scanning for regex '%s'...\n", regex)
		expectedRows = 1
		scanRes = make([]scanTest, 5)
		if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes), document.MatchRegex(regex)); errors.Cause(err) == errors2.ErrDriverFailure {
			fmt.Printf("!! Got anticipated error '%v'\n", err)
		} else if err != nil {
			return errEncountered(err)
		} else if num != expectedRows {
			// there are two entries for this user
			return errEncountered(errors.New("ErrUnexpectedResultLength"))
		} else {
			fmt.Printf("Expected %d results for %s, got %d results:", expectedRows, regex, num)
			for _, row := range scanRes {
				fmt.Printf(" %v", row.User)
			}
			fmt.Println()
		}

		// test Remove
		for k, v := range map[string]string{
			"/test/remove_test_any": `{"Index":"delete_me"}`,
		} {
			if err = json.Unmarshal([]byte(v), &src); err != nil {
				return errEncountered(err)
			}
			fmt.Printf("Setting the value of `%s` to `%s` ...\n", k, v)

			key := document.NewKeyFromStringUnchecked(k)
			if _, err = store.Set(
				key,
				document.WithAnyVersion(),
				document.WithTTL(10*time.Minute),
				document.WithSource(src),
			); err != nil {
				return errEncountered(err)
			}

			// make sure it was set
			if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
				return errEncountered(errors.Wrap(err, "failed to set"))
			} else if !ok {
				return errEncountered(errors.New("ErrKeyNotInStore"))
			}

			// remove it
			fmt.Printf("Removing `%s` ...\n", k)
			if err := store.Remove(key, document.WithAnyVersion()); err != nil {
				return errEncountered(errors.Wrap(err, "failed to remove"))
			}

			// make sure it has been removed
			if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
				return errEncountered(errors.Wrap(err, "failed to verify"))
			} else if ok {
				return errEncountered(errors.New("ErrKeyNotRemoved"))
			}
		}

		// test Remove (with a version provided) - Yes, this is mostly redundant and could be consolidated
		for k, v := range map[string]string{
			"/test/remove_test_ver": `{"Index":"delete_me_too"}`,
		} {
			if err = json.Unmarshal([]byte(v), &src); err != nil {
				return errEncountered(err)
			}
			fmt.Printf("Setting the value of `%s` to `%s` ...\n", k, v)

			key := document.NewKeyFromStringUnchecked(k)
			var version document.Version
			if version, err = store.Set(
				key,
				document.WithAnyVersion(),
				document.WithTTL(10*time.Minute),
				document.WithSource(src),
			); err != nil {
				return errEncountered(err)
			}

			// make sure it was set
			if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
				return errEncountered(errors.Wrap(err, "failed to set"))
			} else if !ok {
				return errEncountered(errors.New("ErrKeyNotInStore"))
			}

			var dst interface{}
			var getVersion document.Version
			if getVersion, err = store.Get(
				key,
				document.WithAnyVersion(),
				document.WithDestination(&dst),
			); err != nil {
				return errEncountered(err)
			} else if version < getVersion {
				/* In some cases (Badger) it is expected for Get() to return a higher version than the immediately
				 * preceding Set() call (because Badger internals). So we only log an warning about this undesirable
				 * case at the time being.
				 */
				fmt.Printf("Warning: Set returned version %d, Get returned verison %d.\n", version, getVersion)
				version = getVersion
			} else if version != getVersion {
				/* However, if the version number ever rolls backward? That's very bad.
				 */
				return errEncountered(errors.Wrap(errors2.ErrInvalidVersioning, "unexpected version rollback"))
			} else {
				fmt.Printf("Got expected version %d back from Get.\n", version)
			}

			// remove it
			fmt.Printf("Removing `%s` version %d ...\n", k, version)
			if err := store.Remove(key, document.WithVersion(version)); err != nil {
				return errEncountered(errors.Wrap(err, "failed to remove"))
			}

			// make sure it has been removed
			if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
				return errEncountered(errors.Wrap(err, "failed to verify"))
			} else if ok {
				return errEncountered(errors.New("ErrKeyNotRemoved"))
			}
		}

	}

	return nil
}

func errEncountered(err error) error {
	fmt.Println("- - - - -")
	fmt.Println(err.Error())
	fmt.Println("- - - - -")
	fmt.Println("		       ┌────────────────────────────────┐")
	fmt.Println("		       │    NoSQL Test Result: FAIL     │")
	fmt.Println("		       └────────────────────────────────┘")
	return err
}
