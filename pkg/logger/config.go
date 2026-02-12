package logger

// Config holds logger settings.
type Config struct {
	Level            string   `yaml:"level"`
	Encoding         string   `yaml:"encoding"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
	Development      bool     `yaml:"development"`
	Sampling         bool     `yaml:"sampling"`
	Stacktrace       string   `yaml:"stacktrace"`
}

// AppInfo identifies the running application instance.
type AppInfo struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"env"`
	Cluster     string `yaml:"cluster"`
	Version     string `yaml:"version"`
	Instance    string `yaml:"instance"`
	Owner       string `yaml:"owner"`
	Process     string `yaml:"process"`
}
