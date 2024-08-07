package pkg

type Cache struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Mount struct {
	Host      string `yaml:"host"`
	Container string `yaml:"container"`
}

type Port struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type EnvVar struct {
	Key       string `yaml:"key"`
	Value     string `yaml:"value"`
	ValueFrom string `yaml:"valueFrom"`
}

type Artifact struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Dependency struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Service struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Env     []EnvVar `yaml:"env"`
	Mounts  []Mount  `yaml:"mounts"`
	Ports   []Port   `yaml:"ports"`
	Command string   `yaml:"command"`
}

type Job struct {
	Name         string            `yaml:"name"`
	Image        string            `yaml:"image"`
	Services     map[string]string `yaml:"services"`
	Caches       []Cache           `yaml:"caches"`
	Mounts       []Mount           `yaml:"mounts"`
	Env          []EnvVar          `yaml:"env"`
	Workdir      string            `yaml:"workdir"`
	Artifacts    []Artifact        `yaml:"artifacts"`
	Dependencies []Dependency      `yaml:"dependencies"`
	Commands     []string          `yaml:"commands"`
}

type SecretEngine struct {
	Name string         `yaml:"name"`
	Type string         `yaml:"type"`
	Spec map[string]any `yaml:"spec"`
}

type Secret struct {
	Name string         `yaml:"name"`
	From string         `yaml:"from"`
	Spec map[string]any `yaml:"spec"`
}

type Stiletto struct {
	SecretEngines []SecretEngine `yaml:"secretEngines"`
	Secrets       []Secret       `yaml:"secrets"`
	Services      []Service      `yaml:"services"`
	Pipeline      []Job          `yaml:"pipeline"`
}
