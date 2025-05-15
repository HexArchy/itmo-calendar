package config

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SecretFilePathSuffix is the suffix appended to flag names for secret file paths.
const SecretFilePathSuffix = "_file_path"

var _envReplacer = strings.NewReplacer(".", "_", "-", "_")

// FlagSet extends pflag.FlagSet to include configuration dictionary and initialization.
type FlagSet struct {
	*pflag.FlagSet
	Dict Dict
	once sync.Once
}

// NewFlagSet creates a new FlagSet with the given name and error handling behavior.
func NewFlagSet(name string, errorHandling pflag.ErrorHandling) *FlagSet {
	fs := pflag.NewFlagSet(name, errorHandling)
	fs.String("config-help", "", "Show help config in spec format: helm/env/yaml")

	return &FlagSet{
		FlagSet: fs,
		Dict:    Dict{},
	}
}

// StateHash returns a hash of the current configuration state.
func (f *FlagSet) StateHash() string {
	st := bytes.Buffer{}
	for _, prop := range f.Dict.Sorted() {
		if prop.Unique {
			continue
		}

		st.WriteString(prop.Name)
		st.WriteString(prop.CurrentVal)
		st.WriteString("\n")
	}

	sum := sha256.Sum256(st.Bytes())

	return hex.EncodeToString(sum[:])
}

// Init initializes the FlagSet with the provided arguments.
func (f *FlagSet) Init(configPath *string, cmdlineArgs ...string) error {
	var err error
	f.once.Do(func() {
		// Register all configurations as flags
		err = viper.BindPFlags(f.FlagSet)
		if err != nil {
			err = errors.Wrap(err, "bind pflags")
			return
		}

		// Parse command line args and fill config with values
		err = f.Parse(cmdlineArgs)
		if err != nil {
			err = errors.Wrap(err, "parse command line arguments")
			return
		}

		// Parse environment variables
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(_envReplacer)

		// Check for config file
		var path string

		// Check if passed via flag
		if configPath != nil {
			path = *configPath
		}

		// Check if passed via ENV
		if path == "" {
			if p, ok := os.LookupEnv(_envConfig); ok {
				path = p
			}
		}

		// If a path was provided in any way
		if path != "" {
			if strings.HasSuffix(path, ".yaml") {
				viper.SetConfigFile(path)
				viper.SetConfigType("yaml")

				// Parse config file
				err = viper.ReadInConfig()
				if err != nil {
					err = errors.Wrap(err, "read configuration file")
					return
				}
			}
			// Other types temporarily unsupported
		}

		// Collect keys that are paths to secret files
		secretMap := make(map[string]string)
		f.VisitAll(func(flag *pflag.Flag) {
			if strings.HasSuffix(flag.Name, SecretFilePathSuffix) && viper.GetString(flag.Name) != "" {
				secretMap[flag.Name] = viper.GetString(flag.Name)
			}
		})

		// Update configurations with additional data
		f.VisitAll(func(flag *pflag.Flag) {
			isSlice := strings.Contains(flag.Value.Type(), "Slice")
			prop := f.Dict.GetOrRegister(flag.Name)
			prop.Description = flag.Usage
			prop.Type = flag.Value.Type()

			if isSlice {
				prop.DefaultVal = sliceDefValue(flag.DefValue)
			} else {
				prop.DefaultVal = flag.DefValue
			}

			// If not set from flag, try to find in viper
			if !flag.Changed {
				if isSlice {
					ss := viper.GetStringSlice(flag.Name)
					sn := make([]string, len(ss))
					for i, v := range ss {
						sn[i] = expandEnvVars(v)
					}
					_ = flag.Value.Set(strings.Join(sn, ","))
					prop.CurrentVal = strings.Join(sn, ",")
				} else {
					val := expandEnvVars(viper.GetString(flag.Name))
					_ = flag.Value.Set(val)
					prop.CurrentVal = flag.Value.String()
				}
			}

			// Check if this is a secret value with a file path
			if path, ok := secretMap[flag.Name+SecretFilePathSuffix]; ok {
				file, fileErr := os.Open(path)
				if fileErr != nil {
					err = errors.Wrap(fileErr, "open secret file")
					return
				}

				rawSecret, readErr := io.ReadAll(file)
				if readErr != nil {
					err = errors.Wrap(readErr, "read secret file")
					_ = file.Close()
					return
				}

				secret := strings.TrimSpace(string(rawSecret))
				_ = flag.Value.Set(secret)

				if closeErr := file.Close(); closeErr != nil {
					err = errors.Wrap(closeErr, "close secret file")
					return
				}
			}

			if isSlice {
				prop.CurrentVal = strings.Join(flag.Value.(pflag.SliceValue).GetSlice(), ",")
			} else {
				prop.CurrentVal = flag.Value.String()
			}
		})

		configHelp, helpErr := f.GetString("config-help")
		if helpErr != nil {
			err = errors.Wrap(helpErr, "get config-help value")
			return
		}

		switch configHelp {
		case "env":
			helpEnv(f.Dict)
		case "helm":
			helpHelm(f.Dict)
		case "yaml":
			helpYaml(f.Dict)
		default:
			if configHelp != "" {
				err = fmt.Errorf("unsupported '%s' format given", configHelp)
				return
			}
		}
	})

	return err
}

// Bool defines a bool flag.
func (f *FlagSet) Bool(name string, defValue bool, description string) *bool {
	return f.FlagSet.Bool(name, defValue, description)
}

// BoolVar defines a bool flag with a specified variable to store the value.
func (f *FlagSet) BoolVar(val *bool, name string, defValue bool, description string) {
	f.FlagSet.BoolVar(val, name, defValue, description)
}

// Int defines an int flag.
func (f *FlagSet) Int(name string, defValue int, description string) *int {
	return f.FlagSet.Int(name, defValue, description)
}

// IntVar defines an int flag with a specified variable to store the value.
func (f *FlagSet) IntVar(val *int, name string, defValue int, description string) {
	f.FlagSet.IntVar(val, name, defValue, description)
}

// Int32 defines an int32 flag.
func (f *FlagSet) Int32(name string, defValue int32, description string) *int32 {
	return f.FlagSet.Int32(name, defValue, description)
}

// Int32Var defines an int32 flag with a specified variable to store the value.
func (f *FlagSet) Int32Var(val *int32, name string, defValue int32, description string) {
	f.FlagSet.Int32Var(val, name, defValue, description)
}

// Int64 defines an int64 flag.
func (f *FlagSet) Int64(name string, defValue int64, description string) *int64 {
	return f.FlagSet.Int64(name, defValue, description)
}

// Int64Var defines an int64 flag with a specified variable to store the value.
func (f *FlagSet) Int64Var(val *int64, name string, defValue int64, description string) {
	f.FlagSet.Int64Var(val, name, defValue, description)
}

// IntSlice defines a []int flag.
func (f *FlagSet) IntSlice(name string, defValue []int, description string) *[]int {
	return f.FlagSet.IntSlice(name, defValue, description)
}

// IntSliceVar defines a []int flag with a specified variable to store the value.
func (f *FlagSet) IntSliceVar(val *[]int, name string, defValue []int, description string) {
	f.FlagSet.IntSliceVar(val, name, defValue, description)
}

// Int32Slice defines a []int32 flag.
func (f *FlagSet) Int32Slice(name string, defValue []int32, description string) *[]int32 {
	return f.FlagSet.Int32Slice(name, defValue, description)
}

// Int32SliceVar defines a []int32 flag with a specified variable to store the value.
func (f *FlagSet) Int32SliceVar(val *[]int32, name string, defValue []int32, description string) {
	f.FlagSet.Int32SliceVar(val, name, defValue, description)
}

// Int64Slice defines a []int64 flag.
func (f *FlagSet) Int64Slice(name string, defValue []int64, description string) *[]int64 {
	return f.FlagSet.Int64Slice(name, defValue, description)
}

// Int64SliceVar defines a []int64 flag with a specified variable to store the value.
func (f *FlagSet) Int64SliceVar(val *[]int64, name string, defValue []int64, description string) {
	f.FlagSet.Int64SliceVar(val, name, defValue, description)
}

// Uint defines a uint flag.
func (f *FlagSet) Uint(name string, defValue uint, description string) *uint {
	return f.FlagSet.Uint(name, defValue, description)
}

// UintVar defines a uint flag with a specified variable to store the value.
func (f *FlagSet) UintVar(val *uint, name string, defValue uint, description string) {
	f.FlagSet.UintVar(val, name, defValue, description)
}

// Uint16 defines a uint16 flag.
func (f *FlagSet) Uint16(name string, defValue uint16, description string) *uint16 {
	return f.FlagSet.Uint16(name, defValue, description)
}

// Uint16Var defines a uint16 flag with a specified variable to store the value.
func (f *FlagSet) Uint16Var(val *uint16, name string, defValue uint16, description string) {
	f.FlagSet.Uint16Var(val, name, defValue, description)
}

// UintSlice defines a []uint flag.
func (f *FlagSet) UintSlice(name string, defValue []uint, description string) *[]uint {
	return f.FlagSet.UintSlice(name, defValue, description)
}

// UintSliceVar defines a []uint flag with a specified variable to store the value.
func (f *FlagSet) UintSliceVar(val *[]uint, name string, defValue []uint, description string) {
	f.FlagSet.UintSliceVar(val, name, defValue, description)
}

// Float32 defines a float32 flag.
func (f *FlagSet) Float32(name string, defValue float32, description string) *float32 {
	return f.FlagSet.Float32(name, defValue, description)
}

// Float32Var defines a float32 flag with a specified variable to store the value.
func (f *FlagSet) Float32Var(val *float32, name string, defValue float32, description string) {
	f.FlagSet.Float32Var(val, name, defValue, description)
}

// Float64 defines a float64 flag.
func (f *FlagSet) Float64(name string, defValue float64, description string) *float64 {
	return f.FlagSet.Float64(name, defValue, description)
}

// Float64Var defines a float64 flag with a specified variable to store the value.
func (f *FlagSet) Float64Var(val *float64, name string, defValue float64, description string) {
	f.FlagSet.Float64Var(val, name, defValue, description)
}

// Float32Slice defines a []float32 flag.
func (f *FlagSet) Float32Slice(name string, defValue []float32, description string) *[]float32 {
	return f.FlagSet.Float32Slice(name, defValue, description)
}

// Float32SliceVar defines a []float32 flag with a specified variable to store the value.
func (f *FlagSet) Float32SliceVar(val *[]float32, name string, defValue []float32, description string) {
	f.FlagSet.Float32SliceVar(val, name, defValue, description)
}

// Float64Slice defines a []float64 flag.
func (f *FlagSet) Float64Slice(name string, defValue []float64, description string) *[]float64 {
	return f.FlagSet.Float64Slice(name, defValue, description)
}

// Float64SliceVar defines a []float64 flag with a specified variable to store the value.
func (f *FlagSet) Float64SliceVar(val *[]float64, name string, defValue []float64, description string) {
	f.FlagSet.Float64SliceVar(val, name, defValue, description)
}

// String defines a string flag.
func (f *FlagSet) String(name string, defValue string, description string) *string {
	return f.FlagSet.String(name, defValue, description)
}

// StringVar defines a string flag with a specified variable to store the value.
func (f *FlagSet) StringVar(val *string, name string, defValue string, description string) {
	f.FlagSet.StringVar(val, name, defValue, description)
}

// StringSlice defines a []string flag.
func (f *FlagSet) StringSlice(name string, defValue []string, description string) *[]string {
	return f.FlagSet.StringSlice(name, defValue, description)
}

// StringSliceVar defines a []string flag with a specified variable to store the value.
func (f *FlagSet) StringSliceVar(val *[]string, name string, defValue []string, description string) {
	f.FlagSet.StringSliceVar(val, name, defValue, description)
}

// Duration defines a time.Duration flag.
func (f *FlagSet) Duration(name string, defValue time.Duration, description string) *time.Duration {
	return f.FlagSet.Duration(name, defValue, description)
}

// DurationVar defines a time.Duration flag with a specified variable to store the value.
func (f *FlagSet) DurationVar(val *time.Duration, name string, defValue time.Duration, description string) {
	f.FlagSet.DurationVar(val, name, defValue, description)
}

// DurationSlice defines a []time.Duration flag.
func (f *FlagSet) DurationSlice(name string, defValue []time.Duration, description string) *[]time.Duration {
	return f.FlagSet.DurationSlice(name, defValue, description)
}

// DurationSliceVar defines a []time.Duration flag with a specified variable to store the value.
func (f *FlagSet) DurationSliceVar(val *[]time.Duration, name string, defValue []time.Duration, description string) {
	f.FlagSet.DurationSliceVar(val, name, defValue, description)
}

// Secret defines a string flag that may be loaded from a file.
func (f *FlagSet) Secret(name string, defValue string, description string) *string {
	// For DT team's flow to get secrets in Kubernetes
	// Currently secrets are passed to Kubernetes in files
	// Register an additional flag with file_path suffix for this purpose
	f.FlagSet.String(name+SecretFilePathSuffix, "", description)

	// Return pointer to the main parameter
	// If a parameter with file_path suffix is provided, the value from the file will be written here
	return f.FlagSet.String(name, defValue, description)
}

// SecretVar defines a string flag that may be loaded from a file with a specified variable.
func (f *FlagSet) SecretVar(val *string, name string, defValue string, description string) {
	// For DT team's flow to get secrets in Kubernetes
	// Currently secrets are passed to Kubernetes in files
	// Register an additional flag with file_path suffix for this purpose
	f.FlagSet.String(name+SecretFilePathSuffix, "", description)

	// Main parameter
	// If a parameter with file_path suffix is provided, the value from the file will be written here
	f.FlagSet.StringVar(val, name, defValue, description)

	// Mark the property as secret.
	f.Dict.GetOrRegister(name).Secret = true
}

// sliceDefValue formats slice default values for display.
func sliceDefValue(val string) string {
	if val == "" {
		return val
	}

	val = strings.Replace(val, "[", "", 1)
	return strings.Replace(val, "]", "", 1)
}

// expandEnvVars replaces ${VAR} patterns in the provided string with
// their corresponding environment variable values
func expandEnvVars(value string) string {
	// Quick check to avoid regex overhead when not needed
	if !strings.Contains(value, "${") {
		return value
	}

	re := regexp.MustCompile(`\${([^}]+)}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name without ${ and }
		varName := match[2 : len(match)-1]

		// Replace with environment variable value
		if val, exists := os.LookupEnv(varName); exists {
			return val
		}

		// Return the original if not found
		return match
	})
}
