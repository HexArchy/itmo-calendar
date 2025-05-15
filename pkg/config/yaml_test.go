package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestYamlFileFromFlag tests loading configuration from a YAML file specified via flag.
func TestYamlFileFromFlag(t *testing.T) {
	// Skip if testdata directory doesn't exist.
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	// Create full path to our test config.
	cfgPath, err := filepath.Abs("testdata/config.yaml")
	require.NoError(t, err)

	// Store original args and manually add flag for config.
	oldArgs := os.Args
	os.Args = append(os.Args, "--config="+cfgPath)

	// Define expected values
	vBool := true
	vInt := 123
	var vInt32 int32 = 123
	var vInt64 int64 = 123
	var vUint uint = 13
	var vUint16 uint16 = 13
	var vFloat32 float32 = 3.14
	vFloat64 := 100.500
	vString := "trololo"
	vDuration := time.Second

	// Define test cases
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "Bool",
			want: &vBool,
			got:  Bool("bool", false, ""),
		},
		{
			name: "Int",
			want: &vInt,
			got:  Int("int", 0, ""),
		},
		{
			name: "Int32",
			want: &vInt32,
			got:  Int32("int32", 0, ""),
		},
		{
			name: "Int64",
			want: &vInt64,
			got:  Int64("int64", 0, ""),
		},
		{
			name: "IntSlice",
			want: &[]int{100, 200, 300},
			got:  IntSlice("int_slice", []int{}, ""),
		},
		{
			name: "Int32Slice",
			want: &[]int32{100, 200, 300},
			got:  Int32Slice("int32_slice", []int32{}, ""),
		},
		{
			name: "Int64Slice",
			want: &[]int64{100, 200, 300},
			got:  Int64Slice("int64_slice", []int64{}, ""),
		},
		{
			name: "Uint",
			want: &vUint,
			got:  Uint("uint", 0, ""),
		},
		{
			name: "Uint16",
			want: &vUint16,
			got:  Uint16("uint16", 0, ""),
		},
		{
			name: "UintSlice",
			want: &[]uint{0, 13, 666},
			got:  UintSlice("uint_slice", []uint{}, ""),
		},
		{
			name: "Float32",
			want: &vFloat32,
			got:  Float32("float32", 0.0, ""),
		},
		{
			name: "Float64",
			want: &vFloat64,
			got:  Float64("float64", 0.0, ""),
		},
		{
			name: "Float32Slice",
			want: &[]float32{3.14, 100.500},
			got:  Float32Slice("float32_slice", []float32{}, ""),
		},
		{
			name: "Float64Slice",
			want: &[]float64{3.14, 100.500},
			got:  Float64Slice("float64_slice", []float64{}, ""),
		},
		{
			name: "String",
			want: &vString,
			got:  String("string", "", ""),
		},
		{
			name: "StringSlice",
			want: &[]string{"foo", "bar"},
			got:  StringSlice("string_slice", []string{}, ""),
		},
		{
			name: "Duration",
			want: &vDuration,
			got:  Duration("duration", time.Millisecond, ""),
		},
		{
			name: "DurationSlice",
			want: &[]time.Duration{time.Second, time.Millisecond, time.Nanosecond},
			got:  DurationSlice("duration_slice", []time.Duration{}, ""),
		},
	}

	// Initialize configs.
	err = InitOnce()
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}

	// Restore original args.
	os.Args = oldArgs
}

// TestYamlFileFromEnvironment tests loading configuration from a YAML file specified via environment variable.
func TestYamlFileFromEnvironment(t *testing.T) {
	// Skip if testdata directory doesn't exist.
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("testEnvironments", pflag.ExitOnError)

	// Create full path to our test config.
	cfgPath, err := filepath.Abs("testdata/config.yaml")
	require.NoError(t, err)

	// Set environment variable with config path.
	err = os.Setenv(_envConfig, cfgPath)
	require.NoError(t, err)

	// Define expected values.
	vBool := true
	vInt := 123
	var vInt32 int32 = 123
	var vInt64 int64 = 123
	var vUint uint = 13
	var vUint16 uint16 = 13
	var vFloat32 float32 = 3.14
	vFloat64 := 100.500
	vString := "trololo"
	vDuration := time.Second

	// Define test cases.
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "Bool",
			want: &vBool,
			got:  f.Bool("bool", false, ""),
		},
		{
			name: "Int",
			want: &vInt,
			got:  f.Int("int", 0, ""),
		},
		{
			name: "Int32",
			want: &vInt32,
			got:  f.Int32("int32", 0, ""),
		},
		{
			name: "Int64",
			want: &vInt64,
			got:  f.Int64("int64", 0, ""),
		},
		{
			name: "IntSlice",
			want: &[]int{100, 200, 300},
			got:  f.IntSlice("int_slice", []int{}, ""),
		},
		{
			name: "Int32Slice",
			want: &[]int32{100, 200, 300},
			got:  f.Int32Slice("int32_slice", []int32{}, ""),
		},
		{
			name: "Int64Slice",
			want: &[]int64{100, 200, 300},
			got:  f.Int64Slice("int64_slice", []int64{}, ""),
		},
		{
			name: "Uint",
			want: &vUint,
			got:  f.Uint("uint", 0, ""),
		},
		{
			name: "Uint16",
			want: &vUint16,
			got:  f.Uint16("uint16", 0, ""),
		},
		{
			name: "UintSlice",
			want: &[]uint{0, 13, 666},
			got:  f.UintSlice("uint_slice", []uint{}, ""),
		},
		{
			name: "Float32",
			want: &vFloat32,
			got:  f.Float32("float32", 0.0, ""),
		},
		{
			name: "Float64",
			want: &vFloat64,
			got:  f.Float64("float64", 0.0, ""),
		},
		{
			name: "Float32Slice",
			want: &[]float32{3.14, 100.500},
			got:  f.Float32Slice("float32_slice", []float32{}, ""),
		},
		{
			name: "Float64Slice",
			want: &[]float64{3.14, 100.500},
			got:  f.Float64Slice("float64_slice", []float64{}, ""),
		},
		{
			name: "String",
			want: &vString,
			got:  f.String("string", "", ""),
		},
		{
			name: "StringSlice",
			want: &[]string{"foo", "bar"},
			got:  f.StringSlice("string_slice", []string{}, ""),
		},
		{
			name: "Duration",
			want: &vDuration,
			got:  f.Duration("duration", time.Millisecond, ""),
		},
		{
			name: "DurationSlice",
			want: &[]time.Duration{time.Second, time.Millisecond, time.Nanosecond},
			got:  f.DurationSlice("duration_slice", []time.Duration{}, ""),
		},
	}

	// Initialize configs.
	err = f.Init(&cfgPath, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}

	// Unset environment variable.
	err = os.Unsetenv(_envConfig)
	require.NoError(t, err)
}

// TestNestingLevelYamlFile tests correct mapping of nested YAML structure to config keys.
func TestNestingLevelYamlFile(t *testing.T) {
	// Skip if testdata directory doesn't exist.
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestNestingLevelYamlFile", pflag.ExitOnError)
	configPath := f.String("config", "", "Configuration file path.")

	// Create full path to our test config.
	cfgPath, err := filepath.Abs("testdata/nesting_level.yaml")
	require.NoError(t, err)

	oldArgs := os.Args
	// Manually add flag with config path.
	os.Args = append(os.Args, "--config="+cfgPath)

	// Define expected values.
	allVars := "test"

	// Define test cases.
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "level1",
			want: &allVars,
			got:  f.String("yaml.level1", "", ""),
		},
		{
			name: "level1WithUnderline",
			want: &allVars,
			got:  f.String("yaml.level1_with_underline", "", ""),
		},
		{
			name: "level2",
			want: &allVars,
			got:  f.String("yaml.level2.value", "", ""),
		},
		{
			name: "level2WithUnderline",
			want: &allVars,
			got:  f.String("yaml.level2.value_with_underline", "", ""),
		},
	}

	// Initialize configs.
	err = f.Init(configPath, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}

	// Restore original args.
	os.Args = oldArgs
}

// TestErrorBadFilePath tests error handling for incorrect file paths.
func TestErrorBadFilePath(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestBadFilePath", pflag.ExitOnError)
	cfgF := "bad_file_path.yaml"

	// Initialize config with expected error.
	err := f.Init(&cfgF, os.Args[1:]...)

	// Run test.
	t.Run("BadFilePath", func(t *testing.T) {
		assert.ErrorContains(t, err, "no such file or directory")
	})
}

// TestPriority tests configuration source priority (flag > env > file > default).
func TestPriority(t *testing.T) {
	// Skip if testdata directory doesn't exist.
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestPriority", pflag.ExitOnError)
	configPath := f.String("config", "", "Configuration file path.")

	// Register all the variable flags BEFORE parsing command line args
	variable1Ptr := f.String("variable1", "default", "")
	variable2Ptr := f.String("variable2", "", "")
	variable3Ptr := f.String("variable3", "", "")
	variable4Ptr := f.String("variable4", "", "")

	// Create full path to our test config.
	cfgPath, err := filepath.Abs("testdata/priority.yaml")
	require.NoError(t, err)

	oldArgs := os.Args
	// Manually add flag with config path.
	os.Args = append(os.Args, "--config="+cfgPath)

	// Set environment variables.
	err = os.Setenv("VARIABLE3", "environment")
	require.NoError(t, err)
	err = os.Setenv("VARIABLE4", "environment")
	require.NoError(t, err)

	// Add flag value (highest priority).
	os.Args = append(os.Args, "--variable4=flag")

	// Define expected values from different sources.
	variable1 := "this_should_be_overridden_by_default" // Default value.
	variable2 := "file"                                 // From YAML file.
	variable3 := "environment"                          // From environment.
	variable4 := "flag"                                 // From command line flag.

	// Initialize configs FIRST
	err = f.Init(configPath, os.Args[1:]...)
	require.NoError(t, err)

	// THEN define test cases and get values AFTER initialization
	tests := []struct {
		name string
		want string
		got  string
	}{
		{
			name: "DefaultValue",
			want: variable1,
			got:  *variable1Ptr,
		},
		{
			name: "FileValue",
			want: variable2,
			got:  *variable2Ptr,
		},
		{
			name: "EnvironmentValue",
			want: variable3,
			got:  *variable3Ptr,
		},
		{
			name: "FlagValue",
			want: variable4,
			got:  *variable4Ptr,
		},
	}

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}

	// Clean up.
	os.Args = oldArgs
	err = os.Unsetenv("VARIABLE3")
	require.NoError(t, err)
	err = os.Unsetenv("VARIABLE4")
	require.NoError(t, err)
}

// TestConfInConf tests environment variable expansion in configuration values.
func TestConfInConf(t *testing.T) {
	// Skip if testdata directory doesn't exist.
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestConfInConf", pflag.ExitOnError)
	configPath := f.String("config", "", "Configuration file path.")

	// Create full path to our test config.
	cfgPath, err := filepath.Abs("testdata/conf_in_conf.yaml")
	require.NoError(t, err)

	oldArgs := os.Args
	// Manually add flag with config path.
	os.Args = append(os.Args, "--config="+cfgPath)
	// Set environment variables.
	err = os.Setenv("BAR", "bar")
	require.NoError(t, err)

	t.Setenv("ENV_VARIABLE", "foo_${BAR}")

	t.Setenv("ENV_VARIABLE_SLICE", "foo_${BAR},${BAR}_FOo")

	// Add flag with environment variable.
	os.Args = append(os.Args, "--flag_variable=foo_${BAR}")

	// Define expected values.
	fileVariable := "foo_bar"
	fileVariableSlice := []string{"foo_bar"}
	envVariable := "foo_bar"
	envVariableSlice := []string{"foo_bar", "bar_FOo"}
	flagVariable := "foo_${BAR}" // Flag values are not expanded.

	// Define test cases.
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "FileValue",
			want: &fileVariable,
			got:  f.String("file_variable", "", ""),
		},
		{
			name: "FileValueSlice",
			want: &fileVariableSlice,
			got:  f.StringSlice("file_variable_slice", []string{}, ""),
		},
		{
			name: "EnvironmentValue",
			want: &envVariable,
			got:  f.String("env_variable", "", ""),
		},
		{
			name: "EnvironmentValueSlice",
			want: &envVariableSlice,
			got:  f.StringSlice("env_variable_slice", []string{}, ""),
		},
		{
			name: "FlagValue",
			want: &flagVariable,
			got:  f.String("flag_variable", "", ""),
		},
	}

	// Initialize configs.
	err = f.Init(configPath, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}

	// Clean up.
	os.Args = oldArgs
	err = os.Unsetenv("BAR")
	require.NoError(t, err)

	err = os.Unsetenv("ENV_VARIABLE")
	require.NoError(t, err)

	err = os.Unsetenv("ENV_VARIABLE_SLICE")
	require.NoError(t, err)
}
