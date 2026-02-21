package config

import "testing"

func TestLoad(t *testing.T) {
	cfg := Load()

	if cfg.Server.Port == 0 {
		t.Error("server port should be set")
	}

	if cfg.Database.Host == "" {
		t.Error("database host should be set")
	}
}
