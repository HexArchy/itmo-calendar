package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigInit tests the public Init function with a struct.
func TestConfigInit(t *testing.T) {
	// Store original state
	oldDefaultFlagSet := _defaultFlagSet
	oldArgs := os.Args

	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestConfigInit", pflag.ExitOnError)

	// Define a configuration struct
	type TestConfig struct {
		ServerName string        `default:"test-server"`
		Port       int           `default:"8080"`
		Timeout    time.Duration `default:"30s"`
		Debug      bool          `default:"true"`
		Nested     struct {
			Value string `default:"nested-value"`
		}
	}

	// Create a new config instance
	cfg := &TestConfig{}

	// Set clean args for this test
	os.Args = append([]string{os.Args[0]}, "--port=9090", "--debug=false")

	// Initialize config
	err := Init(cfg)
	require.NoError(t, err)

	// Verify values
	assert.Equal(t, "test-server", cfg.ServerName)    // Default value
	assert.Equal(t, 9090, cfg.Port)                   // Value from flag
	assert.Equal(t, 30*time.Second, cfg.Timeout)      // Default value
	assert.False(t, cfg.Debug)                        // Value from flag
	assert.Equal(t, "nested-value", cfg.Nested.Value) // Default value

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}

// TestConfigLoadByPrefix tests loading configuration with a prefix.
func TestConfigLoadByPrefix(t *testing.T) {
	// Store original default flag set
	oldDefaultFlagSet := _defaultFlagSet
	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestConfigLoadByPrefix", pflag.ExitOnError)

	// Define a configuration struct
	type TestConfig struct {
		ServerName string `default:"test-server"`
		Port       int    `default:"8080"`
	}

	// Create a new config instance
	cfg := &TestConfig{}

	// Store original args
	oldArgs := os.Args
	// Clear args and add our test args
	os.Args = append([]string{os.Args[0]}, "--app.port=9090")

	// Load config with prefix
	result := LoadByPrefix(cfg, "app")

	// Initialize once to apply the flag values
	err := InitOnce()
	require.NoError(t, err)

	// Verify values
	assert.Equal(t, "test-server", result.ServerName) // Default value
	assert.Equal(t, 9090, result.Port)                // Value from flag with prefix

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}

// TestConfigLoadDefault tests loading default values into a struct.
func TestConfigLoadDefault(t *testing.T) {
	// Store original state
	oldDefaultFlagSet := _defaultFlagSet

	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestConfigLoadDefault", pflag.ExitOnError)

	// Define a configuration struct
	type TestConfig struct {
		ServerName string        `default:"test-server"`
		Port       int           `default:"8080"`
		Timeout    time.Duration `default:"30s"`
	}

	// Create a new config instance
	cfg := &TestConfig{}

	// Load default values only
	result := LoadDefault(cfg)

	// Verify values
	assert.Equal(t, "test-server", result.ServerName)
	assert.Equal(t, 8080, result.Port)
	assert.Equal(t, 30*time.Second, result.Timeout)

	// Restore original state
	_defaultFlagSet = oldDefaultFlagSet
}

// TestCfgStateHash tests the StateHash function.
func TestCfgStateHash(t *testing.T) {
	// Store original default flag set and args
	oldDefaultFlagSet := _defaultFlagSet
	oldArgs := os.Args

	// Create a new flag set for the first test
	_defaultFlagSet = NewFlagSet("TestStateHash1", pflag.ExitOnError)

	// Test with different configs
	_defaultFlagSet.String("test1", "", "Test flag 1")
	_defaultFlagSet.String("test2", "", "Test flag 2")

	os.Args = append([]string{os.Args[0]}, "--test1=value1", "--test2=value2")

	// Initialize default flag set
	err := InitOnce()
	require.NoError(t, err)

	// Get state hash
	hash1 := StateHash()

	// Verify hash is not empty
	assert.NotEmpty(t, hash1)

	// Create a completely new flag set for the second test
	_defaultFlagSet = NewFlagSet("TestStateHash2", pflag.ExitOnError)

	_defaultFlagSet.String("test1", "", "Test flag 1")
	_defaultFlagSet.String("test2", "", "Test flag 2")

	// Change a value and verify hash changes
	os.Args = append([]string{os.Args[0]}, "--test1=value1", "--test2=different")

	// Initialize default flag set again
	err = InitOnce()
	require.NoError(t, err)

	// Get new state hash
	hash2 := StateHash()

	// Verify hash is different
	assert.NotEqual(t, hash1, hash2)

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}

// TestRegisteredProperties tests the RegisteredProperties function.
func TestRegisteredProperties(t *testing.T) {
	// Store original state
	oldDefaultFlagSet := _defaultFlagSet
	oldArgs := os.Args

	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestRegisteredProperties", pflag.ExitOnError)

	_defaultFlagSet.String("prop1", "", "Test property 1")
	_defaultFlagSet.String("prop2", "", "Test property 2")

	// Test with some config properties - reset args to avoid conflicts
	os.Args = append([]string{os.Args[0]}, "--prop1=value1", "--prop2=value2")

	// Initialize default flag set
	err := InitOnce()
	require.NoError(t, err)

	// Get registered properties
	props := RegisteredProperties()

	// Verify properties
	assert.GreaterOrEqual(t, len(props), 2)

	// Create a map for easier lookup
	propMap := make(map[string]*Prop)
	for _, p := range props {
		propMap[p.Name] = p
	}

	// Check for our properties
	assert.Contains(t, propMap, "prop1")
	assert.Contains(t, propMap, "prop2")
	assert.Equal(t, "value1", propMap["prop1"].CurrentVal)
	assert.Equal(t, "value2", propMap["prop2"].CurrentVal)

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}

// TestMarkUnique tests the MarkUnique function.
func TestMarkUnique(t *testing.T) {
	// Store original state
	oldDefaultFlagSet := _defaultFlagSet
	oldArgs := os.Args

	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestMarkUnique", pflag.ExitOnError)

	// Test with a unique property
	os.Args = append([]string{os.Args[0]}, "--unique_prop=unique_value")

	// Define the property
	String("unique_prop", "", "A unique property")

	// Mark it as unique
	MarkUnique("unique_prop")

	// Initialize default flag set
	err := InitOnce()
	require.NoError(t, err)

	// Get registered properties
	props := RegisteredProperties()

	// Find our unique property
	var uniqueProp *Prop
	for _, p := range props {
		if p.Name == "unique_prop" {
			uniqueProp = p
			break
		}
	}

	// Verify property
	require.NotNil(t, uniqueProp)
	assert.True(t, uniqueProp.Unique)
	assert.Equal(t, "unique_value", uniqueProp.CurrentVal)

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}

// TestSecretVar tests the SecretVar function.
func TestSecretVar(t *testing.T) {
	// Create a temporary file for the secret
	tmpFile, err := os.CreateTemp(t.TempDir(), "secret")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write a secret to the file
	_, err = tmpFile.WriteString("super_secret_value")
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	// Store original state
	oldDefaultFlagSet := _defaultFlagSet
	oldArgs := os.Args

	// Create a new flag set for this test
	_defaultFlagSet = NewFlagSet("TestSecretVar", pflag.ExitOnError)

	// Test with a secret property
	os.Args = append([]string{os.Args[0]}, "--secret_prop_file_path="+tmpFile.Name())

	// Define a variable to hold the secret
	var secretValue string

	// Register it as a secret
	SecretVar(&secretValue, "secret_prop", "default_secret", "A secret property")

	// Initialize default flag set
	err = InitOnce()
	require.NoError(t, err)

	// Verify secret value
	assert.Equal(t, "super_secret_value", secretValue)

	// Get registered properties
	props := RegisteredProperties()

	// Find our secret property
	var secretProp *Prop
	for _, p := range props {
		if p.Name == "secret_prop" {
			secretProp = p
			break
		}
	}

	// Verify property
	require.NotNil(t, secretProp)
	assert.True(t, secretProp.Secret)

	// Restore original state
	os.Args = oldArgs
	_defaultFlagSet = oldDefaultFlagSet
}
