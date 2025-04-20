package parser

type Config struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Listener    Listener    `yaml:"listener"`
	Backends    []Backend   `yaml:"backends"`
	HealthCheck HealthCheck `yaml:"health_check"`
}

type Listener struct {
	Protocol string `yaml:"protocol"`
	Port     int    `yaml:"port"`
}

type Backend struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type HealthCheck struct {
	Path               string `yaml:"path"`
	Interval           string `yaml:"interval"`
	Timeout            string `yaml:"timeout"`
	UnhealthyThreshold int    `yaml:"unhealthy_threshold"`
	HealthyThreshold   int    `yaml:"healthy_threshold"`
}
