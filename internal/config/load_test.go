package config

import (
	"path"
	"testing"
)

func Test_Load_NoConfigFile(t *testing.T) {
	missingPath := path.Join(t.TempDir(), "does-not-exist.yaml")

	_, err := Load(missingPath)
	if err == nil {
		t.Fatalf("expected Load to fail for an explicitly passed missing path")
	}

	_, err = Load()
	if err != nil {
		t.Fatalf("Load with no paths should succeed, got: %v", err)
	}
}

func Test_Init_NoConfigFileNextToBinary(t *testing.T) {
	t.Chdir(t.TempDir())

	defaultConfig = Config{}

	_, err := Init()
	if err != nil {
		t.Fatalf("Init should not fail when no config file is present, got: %v", err)
	}
}
