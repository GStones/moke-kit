package document

import (
	"testing"
)

func TestNewKey(t *testing.T) {
	type testCase struct {
		expectedKey   string
		errorExpected bool
	}

	testCases := []testCase{
		// success cases
		{
			expectedKey:   "tests:test.a",
			errorExpected: false,
		},
		{
			expectedKey:   "tests:benchmark:bench",
			errorExpected: false,
		},
		{
			expectedKey:   "tests:benchmark:bench.a",
			errorExpected: false,
		},

		// error cases
		{
			expectedKey:   "",
			errorExpected: true,
		},
		{
			expectedKey:   "t",
			errorExpected: true,
		},
		{
			expectedKey:   "test.testA,testB",
			errorExpected: true,
		},
		{
			expectedKey:   "test/testA/testB",
			errorExpected: true,
		},
	}
	for i, tc := range testCases {
		if key, err := NewKeyFromString(tc.expectedKey); err != nil {
			if tc.errorExpected == true {
				continue
			} else {
				t.Fatal("Unexpected error encountered in test case #", i+1, ":", err)
			}
		} else if key.IsEmpty() || !key.HasValue() {
			t.Fatal("Key in test case #", i+1, "does not have a value.")
		} else if key.String() != tc.expectedKey {
			t.Fatal("Created keys does not match the provided keys in test case #", i+1)
		} else if ok := key.IsEqual(tc.expectedKey); !ok {
			t.Fatal("Created keys's value is not equal to the expected keys in test case #", i+1)
		}
	}
}

func TestNewKeyAndConversions(t *testing.T) {
	if ns := Namespace(); ns != "" {
		t.Errorf("keys namespace (currently `%s`) must not be set during test", ns)
	} else {
		type testCase struct {
			parts         []string
			expectedKey   string
			errorExpected bool
		}

		testCases := []testCase{
			// success cases
			{
				parts:         []string{"a", "b"},
				expectedKey:   "/a/b",
				errorExpected: false,
			},
			{
				parts:         []string{"a", "b", "c"},
				expectedKey:   "/a/b/c",
				errorExpected: false,
			},
			{
				parts:         []string{"a", "b", "c.d"},
				expectedKey:   "/a/b/c.d",
				errorExpected: false,
			},

			// error cases
			{
				parts:         []string{"a"},
				expectedKey:   "",
				errorExpected: true,
			},
			{
				parts:         []string{},
				expectedKey:   "",
				errorExpected: true,
			},
			{
				parts:         []string{"a", "?"},
				expectedKey:   "",
				errorExpected: true,
			},
		}

		for _, tc := range testCases {
			k, e := NewKeyFromParts(tc.parts...)

			if (tc.errorExpected && e == nil) || (!tc.errorExpected && e != nil) {
				t.Errorf("`%v` failed test", tc.parts)
			}

			if e == nil {
				if !k.IsEqual(tc.expectedKey) {
					t.Errorf("`%s` != `%s`", tc.expectedKey, k)
				}

				if !k.IsEqual(string(k.Bytes())) {
					t.Errorf("`%s` != `%s`", tc.expectedKey, k)
				}

				if nk, _ := NewKeyFromParts(k.Parts()...); nk != k {
					t.Errorf("`%s` != `%s`", k, nk)
				}
			}
		}
	}
}

func TestNewKeyFromStringUnchecked(t *testing.T) {
	if ns := Namespace(); ns != "" {
		t.Errorf("Namespace `%s` returned an unexpected value:", ns)
	} else {
		type testCase struct {
			keyString string
		}

		testCases := []testCase{
			{"a, b, c"},
			{"abc"},
			{"abcdefg, hijklmn, opqrstu, vwxyz"},
		}

		for i, tc := range testCases {
			newKey := NewKeyFromStringUnchecked(tc.keyString)
			if newKey.HasValue() == false {
				t.Error("Error: newKey has no value in test case #", i+1)
			}
			if len(newKey.String()) != len(tc.keyString) {
				t.Error("Error: newKey returned an unexpected string length in test case #", i+1)
			}
			if !newKey.IsEqual(newKey.String()) {
				t.Error("Error: newKey is not equal to itself in test case #", i+1)
			}
		}
	}
}

func TestNewKeyFromBytesUnchecked(t *testing.T) {
	if ns := Namespace(); ns != "" {
		t.Errorf("keys namespace (currently `%s`) must not be set during test", ns)
	} else {
		type testCase struct {
			bytesArray []byte
		}

		testCases := []testCase{
			{[]byte("a, b, c")},
			{[]byte("abc")},
			{[]byte("abc, def, ghi")},
		}
		for i, tc := range testCases {
			newKey := NewKeyFromBytesUnchecked(tc.bytesArray)
			if newKey.HasValue() == false {
				t.Error("Error: newKey has no value in test case #", i+1)
			}
			if len(newKey.Bytes()) != len(tc.bytesArray) {
				t.Error("Error: newKey returned an unexpected bytes length in test case #", i+1)
			}
			if !newKey.IsEqual(string(tc.bytesArray)) {
				t.Error("Error: newKey is not equal to itself in test case #", i+1)
			}
		}
	}
}

func TestParseKey(t *testing.T) {
	if ns := Namespace(); ns != "" {
		t.Errorf("keys namespace (currently `%s`) must not be set during test", ns)
	} else {
		type testCase struct {
			key           string
			errorExpected bool
		}

		testCases := []testCase{
			// success cases
			{key: "/a/b", errorExpected: false},
			{key: "/a/b/c", errorExpected: false},
			{key: "/0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_.-/b/c", errorExpected: false},

			// error cases
			{key: "", errorExpected: true},
			{key: "/", errorExpected: true},
			{key: "//", errorExpected: true},
			{key: "/a", errorExpected: true},
			{key: "a/", errorExpected: true},
			{key: "a/b/", errorExpected: true},
			{key: "a/b/c", errorExpected: true},
			{key: "a/b/c/", errorExpected: true},
			{key: `/":?'!@#$%^&*()+~/b`, errorExpected: true},
		}

		for _, tc := range testCases {
			k, e := NewKeyFromString(tc.key)

			if (tc.errorExpected && e == nil) || (!tc.errorExpected && e != nil) {
				t.Errorf("`%s` failed test", tc.key)
			}

			if e == nil && !k.HasValue() {
				t.Errorf("%s doesn't have a value", tc.key)
			}

			if e == nil && !k.IsEqual(tc.key) {
				t.Errorf("`%s` != `%s`", tc.key, k)
			}

			if e == nil && k.IsEmpty() {
				t.Errorf("`%s` was not provided a value`", tc.key)
			}
		}
	}
}

func TestKey_Clear(t *testing.T) {
	type testCase struct {
		expectedKey string
	}

	testCases := []testCase{
		// success cases
		{
			expectedKey: "/a/b",
		},
		{
			expectedKey: "/a/b/c",
		},
	}
	for i, tc := range testCases {
		key := NewKey(tc.expectedKey)
		if key.IsEmpty() || !key.HasValue() {
			t.Fatal("Key in test case #", i+1, "does not have a value.")
		} else if key.String() != tc.expectedKey {
			t.Fatal("Created keys does not match the provided keys in test case #", i+1)
		} else if ok := key.IsEqual(tc.expectedKey); !ok {
			t.Fatal("Created keys's value is not equal to the expected keys in test case #", i+1)
		} else {
			key.Clear()
			if key.value != "" {
				t.Fatal("Key value is not empty in test case #", i+1, "after keys.Clear() call.")
			}
		}
	}
}

func TestKey_Prefixes(t *testing.T) {
	type testCase struct {
		expectedKey    string
		expectedPrefix []string
		errorExpected  bool
	}

	testCases := []testCase{
		// success cases
		{
			expectedKey:    "/a/b",
			expectedPrefix: []string{},
			errorExpected:  false,
		},
		{
			expectedKey:    "/a/0$b",
			expectedPrefix: []string{"0"},
			errorExpected:  false,
		},
		{
			expectedKey:    "/a/0$1$b",
			expectedPrefix: []string{"0", "1"},
			errorExpected:  false,
		},
		{
			expectedKey:    "/0$a/b",
			expectedPrefix: []string{},
			errorExpected:  false,
		},

		// error cases
		{
			expectedKey:   "/a/z$b",
			errorExpected: true,
		},
		{
			expectedKey:   "/a/$$b",
			errorExpected: true,
		},
	}

	for i, tc := range testCases {
		if key, err := NewKeyFromString(tc.expectedKey); err != nil {
			if tc.errorExpected == true {
				continue
			} else {
				t.Fatal("Unexpected error encountered in test case #", i+1, ":", err)
			}
		} else if key.IsEmpty() || !key.HasValue() {
			t.Fatal("Key in test case #", i+1, "does not have a value.")
		} else if key.String() != tc.expectedKey {
			t.Fatal("Created keys does not match the provided keys in test case #", i+1)
		} else if ok := key.IsEqual(tc.expectedKey); !ok {
			t.Fatal("Created keys's value is not equal to the expected keys in test case #", i+1)
		} else if pref := key.Prefixes(); len(pref) != len(tc.expectedPrefix) {
			t.Errorf("Extracted prefixes '%v' do not match expectation of '%v' in test case #%d\n", pref, tc.expectedPrefix, i+1)
		} else {
			for k, v := range pref {
				if tc.expectedPrefix[k] != v {
					t.Errorf("Extracted prefixes '%v' do not match expectation of '%v' in test case #%d\n", pref, tc.expectedPrefix, i+1)
					break
				}
			}
		}
	}
}
