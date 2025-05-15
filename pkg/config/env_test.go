package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnvironments tests retrieving configuration from environment variables.
func TestEnvironments(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("testEnvironments", pflag.ExitOnError)
	cfgFile := ""

	// Set environment variables for testing.
	var err error
	t.Setenv("BOOL", "true")
	t.Setenv("INT", "123")
	t.Setenv("INT32", "123")
	t.Setenv("INT64", "123")
	t.Setenv("INT_SLICE", "100,200,300")
	t.Setenv("INT32_SLICE", "100,200,300")
	t.Setenv("INT64_SLICE", "100,200,300")
	t.Setenv("UINT", "13")
	t.Setenv("UINT16", "13")
	t.Setenv("UINT_SLICE", "0,13,666")
	t.Setenv("FLOAT32", "3.14")
	t.Setenv("FLOAT64", "100.500")
	t.Setenv("FLOAT32_SLICE", "3.14,100.500")
	t.Setenv("FLOAT64_SLICE", "3.14,100.500")
	t.Setenv("STRING", "trololo")
	t.Setenv("STRING_SLICE", "foo,bar")
	t.Setenv("DURATION", "1s")
	t.Setenv("DURATION_SLICE", "1s,1ms,1ns")

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
	err = f.Init(&cfgFile, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// TestNestingLevelEnvironments tests correct mapping between environment variables and config keys.
func TestNestingLevelEnvironments(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestNestingLevelEnvironments", pflag.ExitOnError)
	cfgF := ""

	// Set environment variables for testing.
	t.Setenv("ENV_LEVEL1", "test")
	t.Setenv("ENV_LEVEL1_WITH_UNDERLINE", "test")
	t.Setenv("ENV_LEVEL2_VALUE", "test")
	t.Setenv("ENV_LEVEL2_VALUE_WITH_UNDERLINE", "test")
	t.Setenv("ENV_LEVEL2_WHAT_VALUE_WITH_UNDERLINE", "test")

	// Define expected values
	allVars := "test"

	// Define test cases
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "level1",
			want: &allVars,
			got:  f.String("env.level1", "", ""),
		},
		{
			name: "level1WithUnderline",
			want: &allVars,
			got:  f.String("env.level1_with_underline", "", ""),
		},
		{
			name: "level2",
			want: &allVars,
			got:  f.String("env.level2.value", "", ""),
		},
		{
			name: "level2WithUnderline",
			want: &allVars,
			got:  f.String("env.level2.value_with_underline", "", ""),
		},
		{
			name: "level2WithDash",
			want: &allVars,
			got:  f.String("env.level2-what.value_with_underline", "", ""),
		},
	}

	// Initialize configs.
	err := f.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// TestWrongEnvironmentsData tests handling of incorrect data format in environment variables.
func TestWrongEnvironmentsData(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestWrongEnvironmentsData", pflag.ExitOnError)
	cfgPath := ""

	// Set environment variables with incorrect formats.
	t.Setenv("BOOL", "test")
	t.Setenv("INT", "test")
	t.Setenv("FLOAT", "test")
	t.Setenv("DURATION", "test")
	t.Setenv("INT_SLICE", "1,2,test")

	// Define expected default values.
	var vBool bool
	var vInt int
	var vFloat float64
	var vDuration time.Duration

	// Define test cases.
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		// For incorrect values, defaults should be used.
		{
			name: "Bool",
			want: &vBool,
			got:  f.Bool("bool", true, ""),
		},
		{
			name: "Int",
			want: &vInt,
			got:  f.Int("int", 100500, ""),
		},
		{
			name: "Float",
			want: &vFloat,
			got:  f.Float64("float", 100500, ""),
		},
		{
			name: "Duration",
			want: &vDuration,
			got:  f.Duration("duration", time.Second, ""),
		},
		{
			name: "IntSlice",
			want: &[]int{1, 2, 3},
			got:  f.IntSlice("int_slice", []int{1, 2, 3}, ""),
		},
	}

	// Initialize configs.
	err := f.Init(&cfgPath, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// TestEnvironmentVarsExpansion tests environment variable expansion in config values.
func TestEnvironmentVarsExpansion(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestEnvironmentVarsExpansion", pflag.ExitOnError)
	cfgF := ""

	// Set environment variables.
	t.Setenv("TEST_VAR", "test_value")
	t.Setenv("WITH_ENV_VAR", "prefix_${TEST_VAR}_suffix")

	// Expected expanded value.
	expected := "prefix_test_value_suffix"

	valuePtr := f.String("with_env_var", "", "")

	// Initialize configs.
	err := f.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	// Get the value that should be expanded.
	assert.Equal(t, expected, *valuePtr)
}
