package config

type Logger struct {
	Level            string   `path:"level" default:"info"`
	Encoding         string   `path:"encoding" default:"json"`
	OutputPaths      []string `path:"output_paths" default:"[\"stdout\"]"`
	ErrorOutputPaths []string `path:"error_output_paths" default:"[\"stderr\"]"`
	Development      bool     `path:"development" default:"false"`
	Sampling         bool     `path:"sampling" default:"true"`
	Stacktrace       string   `path:"stacktrace" default:"error"`
}
