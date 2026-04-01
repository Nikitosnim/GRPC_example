package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validConfigYAML = `
env: "local"
storage_path: "./storage/sso.db"
grpc:
  port: 8080
  timeout: "1h"
migrations_path: "/migrations"
token_ttl: "2h"
database:
  host: localhost
  port: 5432
  user: myuser
  password: "123123"
  dbname: mydb
`
)

type mockFetchCfgPathProvider struct {
	path string
}

func (m mockFetchCfgPathProvider) fetchConfigPath() string {
	return m.path
}

func TestMustLoad_Success(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(validConfigYAML), 0644)
	if err != nil {
		t.Fatal("failed to create test config fle:\n" + err.Error())
	}

	// Mock fetch config path
	mockProvider := &mockFetchCfgPathProvider{path: configPath}

	origProvider := cfgPathProvider
	defer func() { cfgPathProvider = origProvider }()
	cfgPathProvider = mockProvider

	// Act
	conf := MustLoad()

	// Assert
	assert.Equal(t, "local", conf.Env)
	assert.Equal(t, "./storage/sso.db", conf.StoragePath)
	assert.Equal(t, 8080, conf.GRPS.Port)
	assert.Equal(t, "/migrations", conf.MigrationsPath)
	assert.Equal(t, "localhost", conf.Db.Host)
	assert.Equal(t, 5432, conf.Db.Port)
	assert.Equal(t, "myuser", conf.Db.User)
	assert.Equal(t, "123123", conf.Db.Password)
	assert.Equal(t, "mydb", conf.Db.Dbname)
}

func TestMustLoad_EmptyPath(t *testing.T) {
	// Mock fetch config path
	mockProvider := &mockFetchCfgPathProvider{path: ""}
	origProvider := cfgPathProvider
	defer func() { cfgPathProvider = origProvider }()
	cfgPathProvider = mockProvider

	// Assert
	assert.PanicsWithValue(t, "config path empty", func() {
		MustLoad()
	}, "Panic is expected when the path is empty.")
}

func TestMustLoad_FileNotFound(t *testing.T) {
	// Mock fetch config path
	nonExistentPath := "/non/ex/path"
	mockProvider := &mockFetchCfgPathProvider{path: nonExistentPath}
	origProvider := cfgPathProvider
	defer func() { cfgPathProvider = origProvider }()
	cfgPathProvider = mockProvider

	// Assert
	assert.Panics(t, func() {
		MustLoad()
	}, "A panic is expected in case of a reading error .yaml file. Not path or invalid file")
}

func TestMustLoad_EnvVar(t *testing.T) {
	tempDir := os.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(validConfigYAML), 0644)
	if err != nil {
		t.Error("failed to create test config fle:\n" + err.Error())
	}

	originalEnv := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalEnv)
	os.Setenv("CONFIG_PATH", configPath)

	conf := MustLoad()

	assert.NotNil(t, conf)
}
