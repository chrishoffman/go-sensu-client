package sensu

import (
	"testing"
)

func Test_SingleConfigFile(t *testing.T) {
	cfg, _ := LoadConfigs("../config/config.json", []string{})

	if cfg.Client.Name != "localhost" {
		t.Errorf("expected %s, got %s", cfg.Client.Name, "localhost")
	}

	if cfg.Client.Address != "192.168.1.1" {
		t.Errorf("expected %s, got %s", cfg.Client.Address, "192.168.1.1")
	}

	if cfg.Client.Subscriptions[0] != "test" {
		t.Errorf("expected %s, got %s", cfg.Client.Subscriptions[0], "test")
	}
}