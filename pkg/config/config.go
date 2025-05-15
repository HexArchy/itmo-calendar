package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var _defaultFlagSet *FlagSet

const _envConfig = "CONFIG"

func init() {
	_defaultFlagSet = NewFlagSet("default", pflag.ExitOnError)
}

// Config interface defines methods for configuration management.
type Config interface {
	// Load initializes configuration from all sources.
	Load() error

	// LoadStruct loads configuration into a struct.
	LoadStruct(interface{}) error

	// LoadByPrefix loads configuration with a prefix.
	LoadByPrefix(interface{}, string) (interface{}, error)

	// GetStateHash returns a hash of the configuration state.
	GetStateHash() string

	// GetRegisteredProperties returns all registered properties.
	GetRegisteredProperties() []*Prop
}

// DefaultConfig is the standard implementation of Config.
type DefaultConfig struct {
	flagSet *FlagSet
}

// New creates a new DefaultConfig instance.
func New() *DefaultConfig {
	return &DefaultConfig{
		flagSet: NewFlagSet("default", pflag.ExitOnError),
	}
}

// StateHash calculates and returns a SHA256 hash of the current configuration state.
// Unique properties are excluded from the hash to allow comparison between instances.
func StateHash() string {
	return _defaultFlagSet.StateHash()
}

// RegisteredProperties returns all registered configuration properties sorted by name.
func RegisteredProperties() []*Prop {
	return _defaultFlagSet.Dict.Sorted()
}

// InitOnce initializes the default configuration once.
func InitOnce() error {
	// Add custom flag for config file
	configPath := _defaultFlagSet.String("config", "", "Configuration file path.")

	return _defaultFlagSet.Init(configPath, os.Args[1:]...)
}

// Init initializes configuration from a struct.
func Init[T any](cfg *T) error {
	loader := NewStructLoader(cfg, _defaultFlagSet)
	err := loader.Load()
	if err != nil {
		return errors.Wrap(err, "load struct")
	}

	return InitOnce()
}

// LoadByPrefix loads configuration into a struct with the given prefix.
func LoadByPrefix[T any](cfg *T, prefix string) *T {
	loader := NewStructLoader(cfg, _defaultFlagSet)
	err := loader.LoadByPrefix(prefix)
	if err != nil {
		panic(err)
	}

	return cfg
}

// LoadDefault loads default values into a struct.
func LoadDefault[T any](cfg *T) *T {
	loader := NewStructLoader(cfg, _defaultFlagSet)
	err := loader.LoadDefault()
	if err != nil {
		panic(err)
	}

	return cfg
}

// Bool registers a boolean flag.
func Bool(name string, defValue bool, description string) *bool {
	return _defaultFlagSet.Bool(name, defValue, description)
}

// BoolVar registers a boolean flag with specified variable.
func BoolVar(val *bool, name string, defValue bool, description string) {
	_defaultFlagSet.BoolVar(val, name, defValue, description)
}

// Int registers an integer flag.
func Int(name string, defValue int, description string) *int {
	return _defaultFlagSet.Int(name, defValue, description)
}

// IntVar registers an integer flag with specified variable.
func IntVar(val *int, name string, defValue int, description string) {
	_defaultFlagSet.IntVar(val, name, defValue, description)
}

// Int32 registers an int32 flag.
func Int32(name string, defValue int32, description string) *int32 {
	return _defaultFlagSet.Int32(name, defValue, description)
}

// Int32Var registers an int32 flag with specified variable.
func Int32Var(val *int32, name string, defValue int32, description string) {
	_defaultFlagSet.Int32Var(val, name, defValue, description)
}

// Int64 registers an int64 flag.
func Int64(name string, defValue int64, description string) *int64 {
	return _defaultFlagSet.Int64(name, defValue, description)
}

// Int64Var registers an int64 flag with specified variable.
func Int64Var(val *int64, name string, defValue int64, description string) {
	_defaultFlagSet.Int64Var(val, name, defValue, description)
}

// IntSlice registers an integer slice flag.
func IntSlice(name string, defValue []int, description string) *[]int {
	return _defaultFlagSet.IntSlice(name, defValue, description)
}

// IntSliceVar registers an integer slice flag with specified variable.
func IntSliceVar(val *[]int, name string, defValue []int, description string) {
	_defaultFlagSet.IntSliceVar(val, name, defValue, description)
}

// Int32Slice registers an int32 slice flag.
func Int32Slice(name string, defValue []int32, description string) *[]int32 {
	return _defaultFlagSet.Int32Slice(name, defValue, description)
}

// Int32SliceVar registers an int32 slice flag with specified variable.
func Int32SliceVar(val *[]int32, name string, defValue []int32, description string) {
	_defaultFlagSet.Int32SliceVar(val, name, defValue, description)
}

// Int64Slice registers an int64 slice flag.
func Int64Slice(name string, defValue []int64, description string) *[]int64 {
	return _defaultFlagSet.Int64Slice(name, defValue, description)
}

// Int64SliceVar registers an int64 slice flag with specified variable.
func Int64SliceVar(val *[]int64, name string, defValue []int64, description string) {
	_defaultFlagSet.Int64SliceVar(val, name, defValue, description)
}

// Uint registers a uint flag.
func Uint(name string, defValue uint, description string) *uint {
	return _defaultFlagSet.Uint(name, defValue, description)
}

// UintVar registers a uint flag with specified variable.
func UintVar(val *uint, name string, defValue uint, description string) {
	_defaultFlagSet.UintVar(val, name, defValue, description)
}

// Uint8 registers a uint8 flag.
func Uint8(name string, defValue uint8, description string) *uint8 {
	return _defaultFlagSet.Uint8(name, defValue, description)
}

// Uint8Var registers a uint8 flag with specified variable.
func Uint8Var(val *uint8, name string, defValue uint8, description string) {
	_defaultFlagSet.Uint8Var(val, name, defValue, description)
}

// Uint16 registers a uint16 flag.
func Uint16(name string, defValue uint16, description string) *uint16 {
	return _defaultFlagSet.Uint16(name, defValue, description)
}

// Uint16Var registers a uint16 flag with specified variable.
func Uint16Var(val *uint16, name string, defValue uint16, description string) {
	_defaultFlagSet.Uint16Var(val, name, defValue, description)
}

// Uint32 registers a uint32 flag.
func Uint32(name string, defValue uint32, description string) *uint32 {
	return _defaultFlagSet.Uint32(name, defValue, description)
}

// Uint32Var registers a uint32 flag with specified variable.
func Uint32Var(val *uint32, name string, defValue uint32, description string) {
	_defaultFlagSet.Uint32Var(val, name, defValue, description)
}

// Uint64 registers a uint64 flag.
func Uint64(name string, defValue uint64, description string) *uint64 {
	return _defaultFlagSet.Uint64(name, defValue, description)
}

// Uint64Var registers a uint64 flag with specified variable.
func Uint64Var(val *uint64, name string, defValue uint64, description string) {
	_defaultFlagSet.Uint64Var(val, name, defValue, description)
}

// UintSlice registers a uint slice flag.
func UintSlice(name string, defValue []uint, description string) *[]uint {
	return _defaultFlagSet.UintSlice(name, defValue, description)
}

// UintSliceVar registers a uint slice flag with specified variable.
func UintSliceVar(val *[]uint, name string, defValue []uint, description string) {
	_defaultFlagSet.UintSliceVar(val, name, defValue, description)
}

// Float32 registers a float32 flag.
func Float32(name string, defValue float32, description string) *float32 {
	return _defaultFlagSet.Float32(name, defValue, description)
}

// Float32Var registers a float32 flag with specified variable.
func Float32Var(val *float32, name string, defValue float32, description string) {
	_defaultFlagSet.Float32Var(val, name, defValue, description)
}

// Float64 registers a float64 flag.
func Float64(name string, defValue float64, description string) *float64 {
	return _defaultFlagSet.Float64(name, defValue, description)
}

// Float64Var registers a float64 flag with specified variable.
func Float64Var(val *float64, name string, defValue float64, description string) {
	_defaultFlagSet.Float64Var(val, name, defValue, description)
}

// Float32Slice registers a float32 slice flag.
func Float32Slice(name string, defValue []float32, description string) *[]float32 {
	return _defaultFlagSet.Float32Slice(name, defValue, description)
}

// Float32SliceVar registers a float32 slice flag with specified variable.
func Float32SliceVar(val *[]float32, name string, defValue []float32, description string) {
	_defaultFlagSet.Float32SliceVar(val, name, defValue, description)
}

// Float64Slice registers a float64 slice flag.
func Float64Slice(name string, defValue []float64, description string) *[]float64 {
	return _defaultFlagSet.Float64Slice(name, defValue, description)
}

// Float64SliceVar registers a float64 slice flag with specified variable.
func Float64SliceVar(val *[]float64, name string, defValue []float64, description string) {
	_defaultFlagSet.Float64SliceVar(val, name, defValue, description)
}

// String registers a string flag.
func String(name string, defValue string, description string) *string {
	return _defaultFlagSet.String(name, defValue, description)
}

// StringVar registers a string flag with specified variable.
func StringVar(val *string, name string, defValue string, description string) {
	_defaultFlagSet.StringVar(val, name, defValue, description)
}

// StringVarP registers a string flag with specified variable and shorthand.
func StringVarP(val *string, name string, shorthand string, defValue string, description string) {
	_defaultFlagSet.StringVarP(val, name, shorthand, defValue, description)
}

// StringSlice registers a string slice flag.
func StringSlice(name string, defValue []string, description string) *[]string {
	return _defaultFlagSet.StringSlice(name, defValue, description)
}

// StringSliceVar registers a string slice flag with specified variable.
func StringSliceVar(val *[]string, name string, defValue []string, description string) {
	_defaultFlagSet.StringSliceVar(val, name, defValue, description)
}

// Duration registers a duration flag.
func Duration(name string, defValue time.Duration, description string) *time.Duration {
	return _defaultFlagSet.Duration(name, defValue, description)
}

// DurationVar registers a duration flag with specified variable.
func DurationVar(val *time.Duration, name string, defValue time.Duration, description string) {
	_defaultFlagSet.DurationVar(val, name, defValue, description)
}

// DurationSlice registers a duration slice flag.
func DurationSlice(name string, defValue []time.Duration, description string) *[]time.Duration {
	return _defaultFlagSet.DurationSlice(name, defValue, description)
}

// DurationSliceVar registers a duration slice flag with specified variable.
func DurationSliceVar(val *[]time.Duration, name string, defValue []time.Duration, description string) {
	_defaultFlagSet.DurationSliceVar(val, name, defValue, description)
}

// Secret registers a secret string flag.
func Secret(name string, defValue string, description string) *string {
	return _defaultFlagSet.Secret(name, defValue, description)
}

// SecretVar registers a secret string flag with specified variable.
func SecretVar(val *string, name string, defValue string, description string) {
	_defaultFlagSet.SecretVar(val, name, defValue, description)
}

// MarkUnique marks a configuration property as unique.
func MarkUnique(name string) {
	_defaultFlagSet.Dict.GetOrRegister(name).Unique = true
}
