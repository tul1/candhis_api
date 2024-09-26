package loadconfig_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
)

type ConfigTestStruct1 struct {
	DBUser     string `yaml:"db_user" validate:"required"`
	DBPassword string `yaml:"db_password" validate:"required"`
	DBHost     string `yaml:"db_host" validate:"required"`
	DBPort     string `yaml:"db_port" validate:"required,numeric"`
	DBName     string `yaml:"db_name" validate:"required"`
}

type ConfigTestStruct2 struct {
	AppName  string `yaml:"app_name" validate:"required"`
	LogLevel string `yaml:"log_level"` // Optional field
	Timeout  int    `yaml:"timeout"`   // Optional field
}

func TestLoadConfigErrorOnNonExistentFile(t *testing.T) {
	var config interface{}
	err := loadconfig.LoadConfig("non_existent.yml", &config)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestLoadConfigErrors(t *testing.T) {
	testCases := map[string]struct {
		configContent string
		expectedError string
	}{
		"Missing Required Field": {
			configContent: `
db_user: "test_user"
db_password: "test_password"
db_host: "localhost"
db_port: "5432"`,
			expectedError: "Key: 'ConfigTestStruct1.DBName' Error:Field validation for 'DBName' failed on the 'required' tag",
		},
		"Invalid Port - Non-Numeric": {
			configContent: `
db_user: "test_user"
db_password: "test_password"
db_host: "localhost"
db_port: "invalid_port"
db_name: "test_db"`,
			expectedError: "Key: 'ConfigTestStruct1.DBPort' Error:Field validation for 'DBPort' failed on the 'numeric' tag",
		},
		"Invalid YAML Format": {
			configContent: `
db_user: "test_user"
db_password: "test_password"
db_host: "localhost"
db_port: 5432
db_name: "test_db`,
			expectedError: "yaml: line 6: found unexpected end of stream",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			var config ConfigTestStruct1
			configFile, err := os.CreateTemp("", "config*.yml")
			assert.NoError(t, err, "Failed to create temp file")
			defer os.Remove(configFile.Name())

			_, err = configFile.WriteString(test.configContent)
			assert.NoError(t, err, "Failed to write to temp file")
			configFile.Close()

			// Load the config
			err = loadconfig.LoadConfig(configFile.Name(), &config)
			assert.ErrorContains(t, err, test.expectedError)
		})
	}
}

func TestLoadConfigSuccessWithOnlyRequiredField(t *testing.T) {
	configContent := `
db_user: "test_user"
db_password: "test_password"
db_host: "localhost"
db_port: "5432"
db_name: "test_db"
`

	configFile, err := os.CreateTemp("", "config1*.yml")
	assert.NoError(t, err, "Failed to create temp file")
	defer os.Remove(configFile.Name())

	_, err = configFile.WriteString(configContent)
	assert.NoError(t, err, "Failed to write to temp file")
	configFile.Close()

	var config1 ConfigTestStruct1
	err = loadconfig.LoadConfig(configFile.Name(), &config1)
	assert.NoError(t, err)

	assert.Equal(t, "test_user", config1.DBUser)
	assert.Equal(t, "test_password", config1.DBPassword)
	assert.Equal(t, "localhost", config1.DBHost)
	assert.Equal(t, "5432", config1.DBPort)
	assert.Equal(t, "test_db", config1.DBName)
}

func TestLoadConfigStruct2(t *testing.T) {
	configContent := `app_name: "TestApp"`

	configFile, err := os.CreateTemp("", "config2*.yml")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(configFile.Name())

	_, err = configFile.WriteString(configContent)
	require.NoError(t, err, "Failed to write to temp file")
	configFile.Close()

	var config2 ConfigTestStruct2
	err = loadconfig.LoadConfig(configFile.Name(), &config2)
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", config2.AppName)
	assert.Equal(t, "", config2.LogLevel)
	assert.Equal(t, 0, config2.Timeout)
}
