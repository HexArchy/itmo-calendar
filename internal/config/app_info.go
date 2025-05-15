package config

type AppInfo struct {
	Name        string `default:"unknown app"`
	Environment string `path:"env" default:"local"`
	Cluster     string `default:"local"`
	Version     string `path:"version"`
	Instance    string `unique:"true"`
	Owner       string `default:"unknown"`
	Process     string `default:"*" desc:"process name"`
}
