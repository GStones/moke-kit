package document

import (
	"testing"
)

func TestNamespace(t *testing.T) {
	type testCase struct {
		testNamespace string
	}
	testCases := []testCase{
		{testNamespace: "test"},
		{testNamespace: "case"},
	}
	for i, tc := range testCases {
		SetNamespace(tc.testNamespace)
		if ns := Namespace(); ns != tc.testNamespace {
			t.Fatal("Namespace was not set to new value in SetNamespace() call for test case #", i+1)
		} else if NamespaceKeyPrefix() != "/"+tc.testNamespace+"/" {
			t.Fatal("Namespace keys prefix was not set to new value in SetNamespace() call for test case #", i+1)
		}
	}
}
