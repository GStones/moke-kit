package utility

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/kelseyhightower/envconfig"
)

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
	if !reflect.ValueOf(spec).Elem().FieldByName("VaultAddr").IsValid() {
		return nil
	}
	conf := &api.Config{
		Address: reflect.ValueOf(spec).Elem().FieldByName("VaultAddr").String(),
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return err
	}
	client.SetToken(reflect.ValueOf(spec).Elem().FieldByName("VaultToken").String())
	secret, err := client.Logical().Read(reflect.ValueOf(spec).Elem().FieldByName("VaultPath").String())
	if err != nil {
		if strings.Contains(err.Error(), "connect: connection refused") ||
			strings.Contains(err.Error(), "connect: connection timed out") {
			log.Println("vault address cannot be connected, only environment configuration available")
			return nil
		} else if strings.Contains(err.Error(), "permission denied") {
			log.Println("vault token is invalid, only environment configuration available")
			return nil
		}
		return err
	}
	m, ok := secret.Data["data"].(map[string]any)
	if !ok {
		return errors.New("can't read from vault")
	}
	t := reflect.TypeOf(spec).Elem()
	num := t.NumField()
	for i := 0; i < num; i++ {
		key := t.Field(i).Tag.Get("vault")
		if key != "" && m[key] != nil {
			reflect.ValueOf(spec).Elem().FieldByName(t.Field(i).Name).SetString(m[key].(string))
		}
	}
	return nil
}
