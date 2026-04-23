package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validConfigYAML = `
env: "local"
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
cache:
  addr: "localhost:6379"
  password: "123123"
  user: "user"
  db: 1
  max_retries: 5
  dial_timeout: "10s"
  timeout: "5s"
  ttl: "300s"
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
	expectedGRPCTimeout, err := time.ParseDuration("1h")
	require.NoError(t, err)

	expectedCacheDialTimeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	expectedCacheTimeout, err := time.ParseDuration("5s")
	require.NoError(t, err)

	expectedCacheTTL, err := time.ParseDuration("300s")
	require.NoError(t, err)

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err = os.WriteFile(configPath, []byte(validConfigYAML), 0644)
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
	assert.Equal(t, 8080, conf.GRPC.Port)
	assert.Equal(t, expectedGRPCTimeout, conf.GRPC.Timeout)
	assert.Equal(t, "/migrations", conf.MigrationsPath)

	// DB config
	assert.Equal(t, "localhost", conf.Db.Host)
	assert.Equal(t, 5432, conf.Db.Port)
	assert.Equal(t, "myuser", conf.Db.User)
	assert.Equal(t, "123123", conf.Db.Password)
	assert.Equal(t, "mydb", conf.Db.Dbname)

	// cache config
	assert.Equal(t, "localhost:6379", conf.Cache.Addr)
	assert.Equal(t, "123123", conf.Cache.Password)
	assert.Equal(t, "user", conf.Cache.User)
	assert.Equal(t, 1, conf.Cache.DB)
	assert.Equal(t, 5, conf.Cache.MaxRetries)
	assert.Equal(t, expectedCacheDialTimeout, conf.Cache.DialTimeout)
	assert.Equal(t, expectedCacheTimeout, conf.Cache.Timeout)
	assert.Equal(t, expectedCacheTTL, conf.Cache.TTL)

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
