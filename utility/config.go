package utility

import (
	"errors"
	"reflect"

	"github.com/hashicorp/vault/api"
	"github.com/kelseyhightower/envconfig"
)

// Load from environment and vault.
// When VaultAddr is configured, Vault read failures are returned (fail-closed).
// When VaultAddr is empty, Vault is skipped and only environment config is used.
func Load(spec any) error {
	if err := envconfig.Process("", spec); err != nil {
		return err
	}
	if err := loadFromVault(spec); err != nil {
		return err
	}
	return nil
}

func loadFromVault(spec any) error {
	rv := reflect.ValueOf(spec)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return nil
	}
	elem := rv.Elem()
	addrField := elem.FieldByName("VaultAddr")
	if !addrField.IsValid() || addrField.Kind() != reflect.String {
		return nil
	}
	vaultAddr := addrField.String()
	if vaultAddr == "" {
		return nil
	}

	tokenField := elem.FieldByName("VaultToken")
	pathField := elem.FieldByName("VaultPath")
	if !tokenField.IsValid() || tokenField.Kind() != reflect.String ||
		!pathField.IsValid() || pathField.Kind() != reflect.String {
		return errors.New("vault config requires VaultToken and VaultPath string fields")
	}
	vaultToken := tokenField.String()
	vaultPath := pathField.String()
	if vaultToken == "" || vaultPath == "" {
		return errors.New("vault token and path are required when VaultAddr is set")
	}

	conf := &api.Config{Address: vaultAddr}
	client, err := api.NewClient(conf)
	if err != nil {
		return err
	}
	client.SetToken(vaultToken)
	secret, err := client.Logical().Read(vaultPath)
	if err != nil {
		return err
	}
	if secret == nil || secret.Data == nil {
		return errors.New("can't read from vault: empty secret")
	}

	data := secret.Data
	if nested, ok := secret.Data["data"].(map[string]any); ok {
		data = nested
	}

	t := elem.Type()
	num := t.NumField()
	for i := 0; i < num; i++ {
		key := t.Field(i).Tag.Get("vault")
		if key == "" || data[key] == nil {
			continue
		}
		field := elem.FieldByName(t.Field(i).Name)
		if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.String {
			continue
		}
		if strVal, ok := data[key].(string); ok {
			field.SetString(strVal)
		}
	}
	return nil
}
