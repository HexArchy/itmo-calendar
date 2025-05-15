package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFlags tests retrieving configuration from command line flags.
func TestFlags(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestFlags", pflag.ExitOnError)
	var cfgF string

	testArgs := append([]string{os.Args[0]},
		"--bool=true",
		"--int=123",
		"--int32=123",
		"--int64=123",
		"--int_slice=100,200,300",
		"--int32_slice=100,200,300",
		"--int64_slice=100,200,300",
		"--uint=13",
		"--uint16=13",
		"--uint_slice=0,13,666",
		"--float32=3.14",
		"--float64=100.500",
		"--float32_slice=3.14,100.500",
		"--float64_slice=3.14,100.500",
		"--string=trololo",
		"--string_slice=foo,bar",
		"--duration=1s",
		"--duration_slice=1s,1ms,1ns")

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

	// Define test cases
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
	err := f.Init(&cfgF, testArgs...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// TestNestingLevelFlags tests correct mapping between flag keys and config keys.
func TestNestingLevelFlags(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestNestingLevelFlags", pflag.ExitOnError)
	var cfgF string

	// Create test arguments
	testArgs := []string{
		"--flag.level1=test",
		"--flag.level1_with_underline=test",
		"--flag.level2.value=test",
		"--flag.level2.value_with_underline=test",
	}

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
			got:  f.String("flag.level1", "", ""),
		},
		{
			name: "level1WithUnderline",
			want: &allVars,
			got:  f.String("flag.level1_with_underline", "", ""),
		},
		{
			name: "level2",
			want: &allVars,
			got:  f.String("flag.level2.value", "", ""),
		},
		{
			name: "level2WithUnderline",
			want: &allVars,
			got:  f.String("flag.level2.value_with_underline", "", ""),
		},
	}

	// Initialize configs.
	err := f.Init(&cfgF, testArgs...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

// TestDefaultValues tests default values when no configs are provided.
func TestDefaultValues(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestDefaultValues", pflag.ExitOnError)
	cfgF := ""

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

	// Define test cases
	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "Bool",
			want: &vBool,
			got:  f.Bool("bool", vBool, ""),
		},
		{
			name: "Int",
			want: &vInt,
			got:  f.Int("int", vInt, ""),
		},
		{
			name: "Int32",
			want: &vInt32,
			got:  f.Int32("int32", vInt32, ""),
		},
		{
			name: "Int64",
			want: &vInt64,
			got:  f.Int64("int64", vInt64, ""),
		},
		{
			name: "IntSlice",
			want: &[]int{100, 200, 300},
			got:  f.IntSlice("int_slice", []int{100, 200, 300}, ""),
		},
		{
			name: "Int32Slice",
			want: &[]int32{100, 200, 300},
			got:  f.Int32Slice("int32_slice", []int32{100, 200, 300}, ""),
		},
		{
			name: "Int64Slice",
			want: &[]int64{100, 200, 300},
			got:  f.Int64Slice("int64_slice", []int64{100, 200, 300}, ""),
		},
		{
			name: "Uint",
			want: &vUint,
			got:  f.Uint("uint", vUint, ""),
		},
		{
			name: "Uint16",
			want: &vUint16,
			got:  f.Uint16("uint16", vUint16, ""),
		},
		{
			name: "UintSlice",
			want: &[]uint{0, 13, 666},
			got:  f.UintSlice("uint_slice", []uint{0, 13, 666}, ""),
		},
		{
			name: "Float32",
			want: &vFloat32,
			got:  f.Float32("float32", vFloat32, ""),
		},
		{
			name: "Float64",
			want: &vFloat64,
			got:  f.Float64("float64", vFloat64, ""),
		},
		{
			name: "Float32Slice",
			want: &[]float32{3.14, 100.500},
			got:  f.Float32Slice("float32_slice", []float32{3.14, 100.500}, ""),
		},
		{
			name: "Float64Slice",
			want: &[]float64{3.14, 100.500},
			got:  f.Float64Slice("float64_slice", []float64{3.14, 100.500}, ""),
		},
		{
			name: "String",
			want: &vString,
			got:  f.String("string", vString, ""),
		},
		{
			name: "StringSlice",
			want: &[]string{"foo", "bar"},
			got:  f.StringSlice("string_slice", []string{"foo", "bar"}, ""),
		},
		{
			name: "Duration",
			want: &vDuration,
			got:  f.Duration("duration", vDuration, ""),
		},
		{
			name: "DurationSlice",
			want: &[]time.Duration{time.Second, time.Millisecond, time.Nanosecond},
			got:  f.DurationSlice("duration_slice", []time.Duration{time.Second, time.Millisecond, time.Nanosecond}, ""),
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

// TestDictionaryDefaults tests that default values are correctly stored in the Dictionary.
func TestDictionaryDefaults(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestDictionaryDefaults", pflag.ExitOnError)
	cfgF := ""

	tests := []struct {
		name        string
		wantDefault string
		got         interface{}
	}{
		{
			name:        "bool_p",
			wantDefault: "true",
			got:         f.Bool("bool_p", true, "bool_p desc"),
		},
		{
			name:        "int_p",
			wantDefault: "100",
			got:         f.Int("int_p", 100, "int_p desc"),
		},
		{
			name:        "int32_p",
			wantDefault: "922",
			got:         f.Int32("int32_p", 922, "int32_p desc"),
		},
		{
			name:        "int64_p",
			wantDefault: "99",
			got:         f.Int64("int64_p", 99, "int64_p desc"),
		},
		{
			name:        "int_slice_p",
			wantDefault: "99,100,101",
			got:         f.IntSlice("int_slice_p", []int{99, 100, 101}, "int_slice_p desc"),
		},
		{
			name:        "int32_slice_p",
			wantDefault: "99,100,102",
			got:         f.Int32Slice("int32_slice_p", []int32{99, 100, 102}, "int32_slice_p desc"),
		},
		{
			name:        "int64_slice_p",
			wantDefault: "100,200,300",
			got:         f.Int64Slice("int64_slice_p", []int64{100, 200, 300}, "int64_slice_p desc"),
		},
		{
			name:        "uint_p",
			wantDefault: "1",
			got:         f.Uint("uint_p", 1, "uint_p desc"),
		},
		{
			name:        "uint16_p",
			wantDefault: "2",
			got:         f.Uint16("uint16_p", 2, "uint16_p desc"),
		},
		{
			name:        "uint_slice_p",
			wantDefault: "0,13,666",
			got:         f.UintSlice("uint_slice_p", []uint{0, 13, 666}, "uint_slice_p desc"),
		},
		{
			name:        "float32_p",
			wantDefault: "4.7",
			got:         f.Float32("float32_p", 4.7, "float32_p desc"),
		},
		{
			name:        "float64_p",
			wantDefault: "444.92",
			got:         f.Float64("float64_p", 444.92, "float64_p desc"),
		},
		{
			name:        "float32_slice_p",
			wantDefault: "3.140000,100.500000",
			got:         f.Float32Slice("float32_slice_p", []float32{3.14, 100.500}, "float32_slice_p desc"),
		},
		{
			name:        "float64_slice_p",
			wantDefault: "3.140000,100.500000",
			got:         f.Float64Slice("float64_slice_p", []float64{3.14, 100.500}, "float64_slice_p desc"),
		},
		{
			name:        "string_p",
			wantDefault: "hello world,ok",
			got:         f.String("string_p", "hello world,ok", "string_p desc"),
		},
		{
			name:        "string_slice_p",
			wantDefault: "foo,bar",
			got:         f.StringSlice("string_slice_p", []string{"foo", "bar"}, "string_slice_p desc"),
		},
		{
			name:        "duration_p",
			wantDefault: "10s",
			got:         f.Duration("duration_p", 10*time.Second, "duration_p desc"),
		},
		{
			name:        "duration_slice_p",
			wantDefault: "1s,1ms,1ns",
			got: f.DurationSlice(
				"duration_slice_p", []time.Duration{
					time.Second, time.Millisecond, time.Nanosecond,
				}, "duration_slice_p desc"),
		},
	}

	// Initialize configs.
	err := f.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := f.Dict[tt.name]
			assert.Equal(t, tt.name, p.Name)
			assert.Equal(t, tt.name+" desc", p.Description)
			assert.Equal(t, tt.wantDefault, p.DefaultVal)
		})
	}
}

// TestDictionaryActualValue tests that actual values are correctly stored in the Dictionary.
func TestDictionaryActualValue(t *testing.T) {
	// Create a new FlagSet since the default one is already used.
	f := NewFlagSet("TestDictionaryActualValue", pflag.ExitOnError)
	cfgF := ""
	testArgs := append([]string{os.Args[0]},
		"--bool_p=false",
		"--int_p=101",
		"--int32_p=123",
		"--int64_p=123",
		"--int_slice_p=100,200,300",
		"--int32_slice_p=100,200,300",
		"--int64_slice_p=100,200,300",
		"--uint_p=13",
		"--uint16_p=13",
		"--uint_slice_p=0,13,666",
		"--float32_p=3.14",
		"--float64_p=100.500",
		"--float32_slice_p=3.14,100.500",
		"--float64_slice_p=3.14,100.500",
		"--string_p=hello world,ok",
		"--string_slice_p=foo,bar",
		"--duration_p=1s",
		"--duration_slice_p=1s,1ms,1ns")

	tests := []struct {
		name string
		want string
		got  interface{}
	}{
		{
			name: "bool_p",
			want: "false",
			got:  f.Bool("bool_p", true, "bool_p desc"),
		},
		{
			name: "int_p",
			want: "101",
			got:  f.Int("int_p", 100, "int_p desc"),
		},
		{
			name: "int32_p",
			want: "123",
			got:  f.Int32("int32_p", 922, "int32_p desc"),
		},
		{
			name: "int64_p",
			want: "123",
			got:  f.Int64("int64_p", 99, "int64_p desc"),
		},
		{
			name: "int_slice_p",
			want: "100,200,300",
			got:  f.IntSlice("int_slice_p", []int{99, 100, 101}, "int_slice_p desc"),
		},
		{
			name: "int32_slice_p",
			want: "100,200,300",
			got:  f.Int32Slice("int32_slice_p", []int32{99, 100, 102}, "int32_slice_p desc"),
		},
		{
			name: "int64_slice_p",
			want: "100,200,300",
			got:  f.Int64Slice("int64_slice_p", []int64{11, 22, 33}, "int64_slice_p desc"),
		},
		{
			name: "uint_p",
			want: "13",
			got:  f.Uint("uint_p", 1, "uint_p desc"),
		},
		{
			name: "uint16_p",
			want: "13",
			got:  f.Uint16("uint16_p", 2, "uint16_p desc"),
		},
		{
			name: "uint_slice_p",
			want: "0,13,666",
			got:  f.UintSlice("uint_slice_p", []uint{0, 13, 12}, "uint_slice_p desc"),
		},
		{
			name: "float32_p",
			want: "3.14",
			got:  f.Float32("float32_p", 4.7, "float32_p desc"),
		},
		{
			name: "float64_p",
			want: "100.5",
			got:  f.Float64("float64_p", 444.92, "float64_p desc"),
		},
		{
			name: "float32_slice_p",
			want: "3.140000,100.500000",
			got:  f.Float32Slice("float32_slice_p", []float32{4, 1.500}, "float32_slice_p desc"),
		},
		{
			name: "float64_slice_p",
			want: "3.140000,100.500000",
			got:  f.Float64Slice("float64_slice_p", []float64{314, 10.500}, "float64_slice_p desc"),
		},
		{
			name: "string_p",
			want: "hello world,ok",
			got:  f.String("string_p", "hello world by default,ok", "string_p desc"),
		},
		{
			name: "string_slice_p",
			want: "foo,bar",
			got:  f.StringSlice("string_slice_p", []string{"s", "b"}, "string_slice_p desc"),
		},
		{
			name: "duration_p",
			want: "1s",
			got:  f.Duration("duration_p", 10*time.Second, "duration_p desc"),
		},
		{
			name: "duration_slice_p",
			want: "1s,1ms,1ns",
			got:  f.DurationSlice("duration_slice_p", []time.Duration{time.Second}, "duration_slice_p desc"),
		},
	}

	// Initialize configs.
	err := f.Init(&cfgF, testArgs...)
	require.NoError(t, err)

	// Run tests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := f.Dict[tt.name]
			assert.Equal(t, tt.name, p.Name)
			assert.Equal(t, tt.want, p.CurrentVal)
			assert.NotEqual(t, p.CurrentVal, p.DefaultVal)
		})
	}
}

// TestStateHash tests that the state hash is correctly calculated.
func TestStateHash(t *testing.T) {
	f1 := NewFlagSet("TestStateHash1", pflag.ExitOnError)
	f2 := NewFlagSet("TestStateHash2", pflag.ExitOnError)

	cfgF := ""

	// Same values in both FlagSets.
	f1.String("string_val", "test", "")
	f1.Int("int_val", 123, "")

	f2.String("string_val", "test", "")
	f2.Int("int_val", 123, "")

	// Initialize configs.
	err := f1.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	err = f2.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	// Hashes should be equal.
	assert.Equal(t, f1.StateHash(), f2.StateHash())

	// Now let's test unique properties.
	f3 := NewFlagSet("TestStateHash3", pflag.ExitOnError)
	f4 := NewFlagSet("TestStateHash4", pflag.ExitOnError)

	f3.String("string_val", "test", "")
	f3.Int("int_val", 123, "")
	f3.String("unique_val", "value1", "")
	f3.Dict.GetOrRegister("unique_val").Unique = true

	f4.String("string_val", "test", "")
	f4.Int("int_val", 123, "")
	f4.String("unique_val", "value2", "")
	f4.Dict.GetOrRegister("unique_val").Unique = true

	err = f3.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	err = f4.Init(&cfgF, os.Args[1:]...)
	require.NoError(t, err)

	// Hashes should still be equal despite different values for unique properties.
	assert.Equal(t, f3.StateHash(), f4.StateHash())
}
