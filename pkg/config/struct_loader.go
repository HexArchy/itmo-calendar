package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
)

var (
	_durationType  = reflect.TypeOf(time.Duration(0))
	_matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	_matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

const (
	_pathTagName        = "path"
	_descriptionTagName = "desc"
	_secretTagName      = "secret"
	_uniqueTagName      = "unique"
)

// StructLoader loads configuration values from a struct.
type StructLoader struct {
	cfg interface{}
	set *FlagSet
}

// NewStructLoader creates a new struct loader instance.
func NewStructLoader(cfg interface{}, set *FlagSet) *StructLoader {
	return &StructLoader{cfg: cfg, set: set}
}

// Load loads configuration values from the struct.
func (l *StructLoader) Load() error {
	return l.LoadByPrefix("")
}

// LoadDefault loads only default values into the struct.
func (l *StructLoader) LoadDefault() error {
	_, err := l.fillDefaults()

	return err
}

// LoadByPrefix loads configuration values with a given prefix.
func (l *StructLoader) LoadByPrefix(prefix string) error {
	t, err := l.fillDefaults()
	if err != nil {
		return errors.Wrap(err, "fill defaults")
	}

	prefix = strings.TrimSpace(prefix)
	var path []string

	if prefix != "" {
		path = []string{prefix}
	} else {
		path = []string{}
	}

	return l.registerProperties(*t, path, "")
}

// fillDefaults initializes the struct with default values.
func (l *StructLoader) fillDefaults() (*reflect.Value, error) {
	t := reflect.ValueOf(l.cfg)
	if t.Kind() != reflect.Pointer {
		return nil, errors.New("pointer is required")
	}

	if t.IsNil() {
		return nil, errors.New("nil is not supported")
	}

	if t.Type().Elem().Kind() != reflect.Struct {
		return nil, errors.New("pointer to struct required")
	}

	initStructNils(reflect.ValueOf(l.cfg))
	err := defaults.Set(l.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "set defaults")
	}

	return &t, nil
}

// registerProperty registers a single property with the flag set.
func (l *StructLoader) registerProperty(path []string, t reflect.Value, tag reflect.StructTag) error {
	name := strings.Join(path, ".")
	prop := &Prop{Name: name}
	defer l.set.Dict.Merge(prop)

	if tval, ok := tag.Lookup(_descriptionTagName); ok {
		prop.Description = tval
	}

	if _, ok := tag.Lookup(_uniqueTagName); ok {
		prop.Unique = true
	}

	if _, ok := tag.Lookup(_secretTagName); ok {
		if t.Kind() != reflect.String {
			return fmt.Errorf("only strings allowed to be a secret, found '%s' in path '%s'",
				t.Kind().String(), prop.Name)
		}

		prop.Secret = true
		l.set.SecretVar(t.Addr().Interface().(*string), prop.Name,
			t.Interface().(string), prop.Description)

		return nil
	}

	if t.Type() == _durationType {
		l.set.DurationVar(t.Addr().Interface().(*time.Duration), prop.Name,
			t.Interface().(time.Duration), prop.Description)

		return nil
	}

	if t.Kind() == reflect.Slice {
		tt := t.Type().Elem()
		if tt.Kind() == reflect.Pointer {
			tt = tt.Elem()
		}

		if tt == _durationType {
			l.set.DurationSliceVar(t.Addr().Interface().(*[]time.Duration), prop.Name,
				t.Interface().([]time.Duration), prop.Description)

			return nil
		}

		switch tt.Kind() {
		case reflect.Int:
			l.set.IntSliceVar(t.Addr().Interface().(*[]int), prop.Name,
				t.Interface().([]int), prop.Description)
		case reflect.Int32:
			l.set.Int32SliceVar(t.Addr().Interface().(*[]int32), prop.Name,
				t.Interface().([]int32), prop.Description)
		case reflect.Int64:
			l.set.Int64SliceVar(t.Addr().Interface().(*[]int64), prop.Name,
				t.Interface().([]int64), prop.Description)
		case reflect.Uint:
			l.set.UintSliceVar(t.Addr().Interface().(*[]uint), prop.Name,
				t.Interface().([]uint), prop.Description)
		case reflect.Float32:
			l.set.Float32SliceVar(t.Addr().Interface().(*[]float32), prop.Name,
				t.Interface().([]float32), prop.Description)
		case reflect.Float64:
			l.set.Float64SliceVar(t.Addr().Interface().(*[]float64), prop.Name,
				t.Interface().([]float64), prop.Description)
		case reflect.String:
			l.set.StringSliceVar(t.Addr().Interface().(*[]string), prop.Name,
				t.Interface().([]string), prop.Description)
		default:
			return fmt.Errorf("unsupported slice kind %s for %s, %s", tt.Kind(), t.Type(), prop.Name)
		}

		return nil
	}

	switch t.Kind() {
	case reflect.Bool:
		l.set.BoolVar(t.Addr().Interface().(*bool), prop.Name,
			t.Interface().(bool), prop.Description)
	case reflect.Int:
		l.set.IntVar(t.Addr().Interface().(*int), prop.Name,
			t.Interface().(int), prop.Description)
	case reflect.Int32:
		l.set.Int32Var(t.Addr().Interface().(*int32), prop.Name,
			t.Interface().(int32), prop.Description)
	case reflect.Int64:
		l.set.Int64Var(t.Addr().Interface().(*int64), prop.Name,
			t.Interface().(int64), prop.Description)
	case reflect.Uint:
		l.set.UintVar(t.Addr().Interface().(*uint), prop.Name,
			t.Interface().(uint), prop.Description)
	case reflect.Uint8:
		l.set.Uint8Var(t.Addr().Interface().(*uint8), prop.Name,
			t.Interface().(uint8), prop.Description)
	case reflect.Uint16:
		l.set.Uint16Var(t.Addr().Interface().(*uint16), prop.Name,
			t.Interface().(uint16), prop.Description)
	case reflect.Uint32:
		l.set.Uint32Var(t.Addr().Interface().(*uint32), prop.Name,
			t.Interface().(uint32), prop.Description)
	case reflect.Uint64:
		l.set.Uint64Var(t.Addr().Interface().(*uint64), prop.Name,
			t.Interface().(uint64), prop.Description)
	case reflect.Float32:
		l.set.Float32Var(t.Addr().Interface().(*float32), prop.Name,
			t.Interface().(float32), prop.Description)
	case reflect.Float64:
		l.set.Float64Var(t.Addr().Interface().(*float64), prop.Name,
			t.Interface().(float64), prop.Description)
	case reflect.String:
		l.set.StringVar(t.Addr().Interface().(*string), prop.Name,
			t.Interface().(string), prop.Description)
	default:
		return fmt.Errorf("unsupported kind %s for %s", t.Kind(), t.Type())
	}

	return nil
}

// registerProperties recursively registers all properties in a struct.
func (l *StructLoader) registerProperties(t reflect.Value, path []string, tag reflect.StructTag) error {
	if t.Type() == _durationType {
		return l.registerProperty(path, t, tag)
	}

	switch t.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.String,
		reflect.Float64,
		reflect.Slice:
		return l.registerProperty(path, t, tag)
	case reflect.Pointer,
		reflect.UnsafePointer:
		return l.registerProperties(t.Elem(), path, tag)
	case reflect.Struct:
		sf := structFields(t)

		for _, v := range sf {
			name := toSnakeCase(v.Name)

			if tval, ok := v.Tag.Lookup(_pathTagName); ok {
				name = tval
			}

			if name == "-" {
				continue
			}

			err := l.registerProperties(t.FieldByName(v.Name), append(path, name), v.Tag)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("register property %s", v.Name))
			}
		}

		return nil

	default:
		return fmt.Errorf("unsupported kind %s", t.Kind().String())
	}
}

// initStructNils recursively initializes nil pointers in a struct.
func initStructNils(t reflect.Value) {
	if t.Kind() == reflect.Pointer || t.Kind() == reflect.UnsafePointer {
		if t.IsNil() {
			newValue := reflect.New(t.Type().Elem())
			t.Set(newValue)
			initStructNils(newValue)

			return
		} else {
			initStructNils(t.Elem())

			return
		}
	}

	if t.Kind() == reflect.Struct {
		sf := structFields(t)

		for _, v := range sf {
			initStructNils(t.FieldByName(v.Name))
		}
	}
}

// structFields returns all exported fields of a struct.
func structFields(value reflect.Value) []reflect.StructField {
	t := value.Type()
	var f []reflect.StructField

	for i := range t.NumField() {
		field := t.Field(i)

		// We can't access the value of unexported fields.
		if field.PkgPath != "" {
			continue
		}

		f = append(f, field)
	}

	return f
}

// toSnakeCase converts a camelCase string to snake_case.
func toSnakeCase(str string) string {
	snake := _matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = _matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}
