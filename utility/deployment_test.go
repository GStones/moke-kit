package utility

import "testing"

func TestParseDeploymentsNormalization(t *testing.T) {
	if got := ParseDeployments(" PROD "); got != DeploymentsProd {
		t.Fatalf("expected prod, got %q", got)
	}
	if got := ParseDeployments("Dev"); got != DeploymentsDev {
		t.Fatalf("expected dev, got %q", got)
	}
}

func TestDeploymentMatchers(t *testing.T) {
	cases := []struct {
		value            string
		prod, dev, local bool
	}{
		{"prod", true, false, false},
		{"prod_gcp", true, false, false},
		{"prod-aws", true, false, false},
		{"production", false, false, false},
		{"notprod", false, false, false},
		{"nonprod", false, false, false},
		{"dev", false, true, false},
		{"dev_test", false, true, false},
		{"development", false, false, false},
		{"local", false, false, true},
		{"local_67", false, false, true},
		{"localprod", false, false, false},
	}
	for _, tc := range cases {
		d := ParseDeployments(tc.value)
		if d.IsProd() != tc.prod || d.IsDev() != tc.dev || d.IsLocal() != tc.local {
			t.Fatalf("%q: got prod=%v dev=%v local=%v, want prod=%v dev=%v local=%v",
				tc.value, d.IsProd(), d.IsDev(), d.IsLocal(), tc.prod, tc.dev, tc.local)
		}
	}
}
