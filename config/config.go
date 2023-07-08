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
	Label    string     `yaml:"label"`
	Protocol []Protocol `yaml:"protocol"`
}

type Protocol struct {
	Name       string  `yaml:"name"`
	Help       string  `yaml:"help"`
	Label      string  `yaml:"label"`
	Datatype   string  `yaml:"datatype"`
	TrueValue  int     `yaml:"trueValue"`
	MetricType string  `yaml:"metricType"`
	Offset     float64 `yaml:"offset"`
}
