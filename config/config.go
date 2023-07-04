package config

type Config struct {
	Listening Listening  `yaml:"listening"`
	Endpoint  []Endpoint `yaml:"endpoint"`
}

type Listening struct {
	Address string `yaml:"address"`
	Path    string `yaml:"path"`
}

type Endpoint struct {
	Address  string     `yaml:"address"`
	Type     string     `yaml:"type"`
	Length   int        `yaml:"length"`
	Inorder  bool       `yaml:"inorder"`
	Label    string     `yaml:"label"`
	Protocol []Protocol `yaml:"protocol"`
}

type Protocol struct {
	Name       string  `yaml:"name"`
	Help       string  `yaml:"help"`
	Datatype   string  `yaml:"datatype"`
	MetricType string  `yaml:"metricType"`
	Offset     float64 `yaml:"offset"`
}
