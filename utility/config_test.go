package utility

import "testing"

func TestLoadSkipsVaultWhenAddrEmpty(t *testing.T) {
	type cfg struct {
		Name       string `envconfig:"TEST_CFG_NAME" default:"ok"`
		VaultAddr  string `envconfig:"TEST_CFG_VAULT_ADDR" default:""`
		VaultToken string `envconfig:"TEST_CFG_VAULT_TOKEN" default:""`
		VaultPath  string `envconfig:"TEST_CFG_VAULT_PATH" default:""`
	}
	var c cfg
	if err := Load(&c); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if c.Name != "ok" {
		t.Fatalf("expected default name, got %q", c.Name)
	}
}

func TestLoadFailsClosedWhenVaultConfiguredButIncomplete(t *testing.T) {
	type cfg struct {
		VaultAddr  string `envconfig:"TEST_CFG2_VAULT_ADDR" default:"http://127.0.0.1:8200"`
		VaultToken string `envconfig:"TEST_CFG2_VAULT_TOKEN" default:""`
		VaultPath  string `envconfig:"TEST_CFG2_VAULT_PATH" default:""`
	}
	var c cfg
	if err := Load(&c); err == nil {
		t.Fatal("expected error when VaultAddr set without token/path")
	}
}
