package config

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructLoader(t *testing.T) {
	t.Run("should not accept unexpected values", func(t *testing.T) {
		f := NewFlagSet("TestStructLoader", pflag.ExitOnError)
		cfg := ""

		// Initialize configs
		err := f.Init(&cfg, os.Args[1:]...)
		require.NoError(t, err)

		tests := []interface{}{
			nil,
			map[string]interface{}{},
			[]string{},
			1,
			false,
			484.3,
			struct {
				A string
			}{"hello"},
		}

		for _, tt := range tests {
			var des string
			if tt != nil {
				des = reflect.TypeOf(tt).String()
			} else {
				des = "nil"
			}

			t.Run(des, func(t *testing.T) {
				le := len(f.Dict)

				loader := NewStructLoader(tt, f)
				err := loader.Load()
				require.Error(t, err)
				assert.Len(t, f.Dict, le)
			})
		}
	})

	t.Run("should register all types with description", func(t *testing.T) {
		f := NewFlagSet("StructLoaderValidationTest", pflag.ExitOnError)

		config := struct {
			BoolP          bool            `desc:"bool_p desc"`
			IntP           int             `desc:"int_p desc"`
			Int32P         int32           `desc:"int32_p desc"`
			Int64P         int64           `desc:"int64_p desc"`
			IntSliceP      []int           `desc:"int_slice_p desc"`
			Int32SliceP    []int32         `desc:"int32_slice_p desc"`
			Int64SliceP    []int64         `desc:"int64_slice_p desc"`
			UintP          uint            `desc:"uint_p desc"`
			Uint8P         uint8           `desc:"uint8_p desc"`
			Uint16P        uint16          `desc:"uint16_p desc"`
			Uint32P        uint32          `desc:"uint32_p desc"`
			Uint64P        uint64          `desc:"uint64_p desc"`
			UintSliceP     []uint          `desc:"uint_slice_p desc"`
			Float32P       float32         `desc:"float32_p desc"`
			Float64P       float64         `desc:"float64_p desc"`
			Float32SliceP  []float32       `desc:"float32_slice_p desc"`
			Float64SliceP  []float64       `desc:"float64_slice_p desc"`
			StringP        string          `desc:"string_p desc"`
			StringSliceP   []string        `desc:"string_slice_p desc"`
			DurationP      time.Duration   `desc:"duration_p desc"`
			DurationSliceP []time.Duration `desc:"duration_slice_p desc"`
			Level2         *struct {
				Prop1         string `desc:"level2.prop1 desc"`
				Level3LikeRPC *struct {
					Prop222World int `desc:"level2.level3_like_rpc.prop222_world desc"`
				}
			}
		}{}

		loader := NewStructLoader(&config, f)
		err := loader.Load()
		require.NoError(t, err)

		cfg := ""
		err = f.Init(&cfg, os.Args[1:]...)
		require.NoError(t, err)

		tests := []string{
			"bool_p",
			"int_p",
			"int32_p",
			"int64_p",
			"int_slice_p",
			"int32_slice_p",
			"int64_slice_p",
			"uint_p",
			"uint8_p",
			"uint16_p",
			"uint32_p",
			"uint64_p",
			"uint_slice_p",
			"float32_p",
			"float64_p",
			"float32_slice_p",
			"float64_slice_p",
			"string_p",
			"string_slice_p",
			"duration_p",
			"duration_slice_p",
			"level2.prop1",
			"level2.level3_like_rpc.prop222_world",
		}

		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				p, ok := f.Dict[tt]
				require.True(t, ok)
				assert.Equal(t, tt+" desc", p.Description)
			})
		}
	})

	t.Run("should register defaults for all types", func(t *testing.T) {
		f := NewFlagSet("StructLoaderValidationTest1", pflag.ExitOnError)

		config := struct {
			BoolP         bool          `default:"true"`
			IntP          int           `default:"1"`
			Int32P        int32         `default:"2"`
			Int64P        int64         `default:"3"`
			IntSliceP     []int         `default:"[1,2,3]"`
			Int32SliceP   []int32       `default:"[4,5,6]"`
			Int64SliceP   []int64       `default:"[7,8,9]"`
			UintP         uint          `default:"1"`
			Uint8P        uint8         `default:"2"`
			Uint16P       uint16        `default:"3"`
			Uint32P       uint32        `default:"4"`
			Uint64P       uint64        `default:"5"`
			UintSliceP    []uint        `default:"[1,2,3,4,5]"`
			Float32P      float32       `default:"1.1"`
			Float64P      float64       `default:"1.2"`
			Float32SliceP []float32     `default:"[1.3,1.4]"`
			Float64SliceP []float64     `default:"[1.5,1.6]"`
			StringP       string        `default:"hello world"`
			StringSliceP  []string      `default:"[\"again\",\"hello\"]"`
			DurationP     time.Duration `default:"1m17s"`
			Level2        *struct {
				Prop1         string `default:"hello?"`
				Level3LikeRPC *struct {
					Prop222World int `default:"10"`
				}
			}
		}{}

		loader := NewStructLoader(&config, f)
		err := loader.Load()
		require.NoError(t, err)

		cfg := ""
		err = f.Init(&cfg, os.Args[1:]...)
		require.NoError(t, err)

		tests := []struct {
			name string
			want string
		}{
			{
				name: "bool_p",
				want: "true",
			},
			{
				name: "int_p",
				want: "1",
			},
			{
				name: "int32_p",
				want: "2",
			},
			{
				name: "int64_p",
				want: "3",
			},
			{
				name: "int_slice_p",
				want: "1,2,3",
			},
			{
				name: "int32_slice_p",
				want: "4,5,6",
			},
			{
				name: "int64_slice_p",
				want: "7,8,9",
			},
			{
				name: "uint_p",
				want: "1",
			},
			{
				name: "uint8_p",
				want: "2",
			},
			{
				name: "uint16_p",
				want: "3",
			},
			{
				name: "uint32_p",
				want: "4",
			},
			{
				name: "uint64_p",
				want: "5",
			},
			{
				name: "uint_slice_p",
				want: "1,2,3,4,5",
			},
			{
				name: "float32_p",
				want: "1.1",
			},
			{
				name: "float64_p",
				want: "1.2",
			},
			{
				name: "float32_slice_p",
				want: "1.300000,1.400000",
			},
			{
				name: "float64_slice_p",
				want: "1.500000,1.600000",
			},
			{
				name: "string_p",
				want: "hello world",
			},
			{
				name: "string_slice_p",
				want: "again,hello",
			},
			{
				name: "duration_p",
				want: "1m17s",
			},
			{
				name: "level2.prop1",
				want: "hello?",
			},
			{
				name: "level2.level3_like_rpc.prop222_world",
				want: "10",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, ok := f.Dict[tt.name]
				require.True(t, ok)
				assert.Equal(t, tt.want, p.DefaultVal)
			})
		}
	})

	t.Run("should fill struct for all types", func(t *testing.T) {
		f := NewFlagSet("StructLoaderValidationTest1", pflag.ExitOnError)

		config := struct {
			BoolP          bool          `default:"true"`
			IntP           int           `default:"1"`
			Int32P         int32         `default:"2"`
			Int64P         int64         `default:"3"`
			IntSliceP      []int         `default:"[1,2,3]"`
			Int32SliceP    []int32       `default:"[4,5,6]"`
			Int64SliceP    []int64       `default:"[7,8,9]"`
			UintP          uint          `default:"1"`
			Uint8P         uint8         `default:"2"`
			Uint16P        uint16        `default:"3"`
			Uint32P        uint32        `default:"4"`
			Uint64P        uint64        `default:"5"`
			UintSliceP     []uint        `default:"[1,2,3,4,5]"`
			Float32P       float32       `default:"1.1"`
			Float64P       float64       `default:"1.2"`
			Float32SliceP  []float32     `default:"[1.3,1.4]"`
			Float64SliceP  []float64     `default:"[1.5,1.6]"`
			StringP        string        `default:"hello world"`
			StringSliceP   []string      `default:"[\"again\",\"hello\"]"`
			DurationP      time.Duration `default:"1m17s"`
			DurationSliceP []time.Duration
			Level2         *struct {
				Prop1         string `default:"hello?"`
				Level3LikeRPC *struct {
					Prop222World int `default:"10"`
				}
			}
		}{}

		loader := NewStructLoader(&config, f)
		err := loader.Load()
		require.NoError(t, err)

		oldArgs := os.Args
		os.Args = append(os.Args, "--bool_p=false")
		os.Args = append(os.Args, "--int_p=-1")
		os.Args = append(os.Args, "--int32_p=1")
		os.Args = append(os.Args, "--int64_p=2")
		os.Args = append(os.Args, "--int_slice_p=3,-4")
		os.Args = append(os.Args, "--int32_slice_p=100,200,500")
		os.Args = append(os.Args, "--int64_slice_p=-100,-200")
		os.Args = append(os.Args, "--uint_p=2")
		os.Args = append(os.Args, "--uint8_p=3")
		os.Args = append(os.Args, "--uint16_p=4")
		os.Args = append(os.Args, "--uint32_p=5")
		os.Args = append(os.Args, "--uint64_p=6")
		os.Args = append(os.Args, "--uint_slice_p=32,92")
		os.Args = append(os.Args, "--float32_p=32.92")
		os.Args = append(os.Args, "--float32_slice_p=32.2,-1.1,0.1")
		os.Args = append(os.Args, "--float64_p=1.0001")
		os.Args = append(os.Args, "--float64_slice_p=0.0007,-10.01")
		os.Args = append(os.Args, "--string_p=hello,world,like")
		os.Args = append(os.Args, "--string_slice_p=hello,world,like")
		os.Args = append(os.Args, "--duration_p=10m")
		os.Args = append(os.Args, "--duration_slice_p=10m,7m,10s")
		os.Args = append(os.Args, "--level2.prop1=ooooopsddd")
		os.Args = append(os.Args, "--level2.level3_like_rpc.prop222_world=8789889")

		cfg := ""
		err = f.Init(&cfg, os.Args[1:]...)
		require.NoError(t, err)

		tests := []struct {
			name string
			want interface{}
			get  func() interface{}
		}{
			{
				name: "bool_p",
				want: false,
				get: func() interface{} {
					return config.BoolP
				},
			},
			{
				name: "int_p",
				want: -1,
				get: func() interface{} {
					return config.IntP
				},
			},
			{
				name: "int32_p",
				want: int32(1),
				get: func() interface{} {
					return config.Int32P
				},
			},
			{
				name: "int64_p",
				want: int64(2),
				get: func() interface{} {
					return config.Int64P
				},
			},
			{
				name: "int_slice_p",
				want: []int{3, -4},
				get: func() interface{} {
					return config.IntSliceP
				},
			},
			{
				name: "int32_slice_p",
				want: []int32{100, 200, 500},
				get: func() interface{} {
					return config.Int32SliceP
				},
			},
			{
				name: "int64_slice_p",
				want: []int64{-100, -200},
				get: func() interface{} {
					return config.Int64SliceP
				},
			},
			{
				name: "uint_p",
				want: uint(2),
				get: func() interface{} {
					return config.UintP
				},
			},
			{
				name: "uint8_p",
				want: uint8(3),
				get: func() interface{} {
					return config.Uint8P
				},
			},
			{
				name: "uint16_p",
				want: uint16(4),
				get: func() interface{} {
					return config.Uint16P
				},
			},
			{
				name: "uint32_p",
				want: uint32(5),
				get: func() interface{} {
					return config.Uint32P
				},
			},
			{
				name: "uint64_p",
				want: uint64(6),
				get: func() interface{} {
					return config.Uint64P
				},
			},
			{
				name: "uint_slice_p",
				want: []uint{32, 92},
				get: func() interface{} {
					return config.UintSliceP
				},
			},
			{
				name: "float32_p",
				want: float32(32.92),
				get: func() interface{} {
					return config.Float32P
				},
			},
			{
				name: "float64_p",
				want: 1.0001,
				get: func() interface{} {
					return config.Float64P
				},
			},
			{
				name: "float32_slice_p",
				want: []float32{32.2, -1.1, 0.1},
				get: func() interface{} {
					return config.Float32SliceP
				},
			},
			{
				name: "float64_slice_p",
				want: []float64{0.0007, -10.01},
				get: func() interface{} {
					return config.Float64SliceP
				},
			},
			{
				name: "string_p",
				want: "hello,world,like",
				get: func() interface{} {
					return config.StringP
				},
			},
			{
				name: "string_slice_p",
				want: []string{"hello", "world", "like"},
				get: func() interface{} {
					return config.StringSliceP
				},
			},
			{
				name: "duration_p",
				want: 10 * time.Minute,
				get: func() interface{} {
					return config.DurationP
				},
			},
			{
				name: "duration_slice_p",
				want: []time.Duration{
					10 * time.Minute,
					7 * time.Minute,
					10 * time.Second,
				},
				get: func() interface{} {
					return config.DurationSliceP
				},
			},
			{
				name: "level2.prop1",
				want: "ooooopsddd",
				get: func() interface{} {
					return config.Level2.Prop1
				},
			},
			{
				name: "level2.level3_like_rpc.prop222_world",
				want: 8789889,
				get: func() interface{} {
					return config.Level2.Level3LikeRPC.Prop222World
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.EqualValues(t, tt.want, tt.get())
			})
		}

		os.Args = oldArgs
	})

	t.Run("should override default for nested struct", func(t *testing.T) {
		f := NewFlagSet("DefaultOverrideForNestedTest", pflag.ExitOnError)

		config := struct {
			Nested struct {
				Prop1 int `default:"1"`
				Prop2 int `default:"2"`
			} `default:"{\"prop1\":111}"`
		}{}

		loader := NewStructLoader(&config, f)
		err := loader.Load()
		require.NoError(t, err)

		assert.Equal(t, 2, config.Nested.Prop2)
		assert.Equal(t, 111, config.Nested.Prop1)
	})

	t.Run("should work with path override", func(t *testing.T) {
		f := NewFlagSet("StructLoaderValidationTest", pflag.ExitOnError)

		config := struct {
			IntP    int    `desc:"int_p desc"`
			StringP string `path:"striiiing" desc:"striiiing desc"`
			Level2  *struct {
				Prop1         string `desc:"nested_level.prop1 desc"`
				Level3LikeRPC *struct {
					Prop222World int `desc:"nested_level.Level3LikeRPC.prop222_world desc"`
				} `path:"Level3LikeRPC"`
			} `path:"nested_level"`
		}{}

		loader := NewStructLoader(&config, f)
		err := loader.Load()
		require.NoError(t, err)

		cfg := ""
		err = f.Init(&cfg, os.Args[1:]...)
		require.NoError(t, err)

		tests := []string{
			"int_p",
			"striiiing",
			"nested_level.prop1",
			"nested_level.Level3LikeRPC.prop222_world",
		}

		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				p, ok := f.Dict[tt]
				require.True(t, ok)
				assert.Equal(t, tt+" desc", p.Description)
			})
		}
	})
}

// TestToSnakeCase tests converting camelCase to snake_case.
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "helloWorld",
			want:  "hello_world",
		},
		{
			input: "HelloWorld",
			want:  "hello_world",
		},
		{
			input: "HTTPRequest",
			want:  "http_request",
		},
		{
			input: "URLParser",
			want:  "url_parser",
		},
		{
			input: "ID",
			want:  "id",
		},
		{
			input: "UserID",
			want:  "user_id",
		},
		{
			input: "UserUrlAPI",
			want:  "user_url_api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestStructFields tests extracting exported fields from a struct.
func TestStructFields(t *testing.T) {
	type testStruct struct {
		ExportedField   string
		unexportedField string // This should be ignored.
		OtherField      int
	}

	ts := testStruct{
		ExportedField:   "value1",
		unexportedField: "value2",
		OtherField:      123,
	}

	fields := structFields(reflect.ValueOf(ts))

	assert.Len(t, fields, 2) // Only the 2 exported fields.

	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}

	assert.Contains(t, names, "ExportedField")
	assert.Contains(t, names, "OtherField")
	assert.NotContains(t, names, "unexportedField")
}

// TestInitStructNils tests initializing nil pointers in a struct.
func TestInitStructNils(t *testing.T) {
	type deepNested struct {
		Value int
	}

	type nested struct {
		Deep *deepNested
	}

	type testStruct struct {
		Nested *nested
	}

	// Create a struct with nil pointers.
	var ts testStruct

	// Initialize nil pointers.
	initStructNils(reflect.ValueOf(&ts))

	// Check that nil pointers are now initialized.
	assert.NotNil(t, ts.Nested)
	assert.NotNil(t, ts.Nested.Deep)
}
