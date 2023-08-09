package pkg

type Cache struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Mount struct {
	Host      string `yaml:"host"`
	Container string `yaml:"container"`
}

type Job struct {
	Name     string            `yaml:"name"`
	Image    string            `yaml:"image"`
	Caches   []Cache           `yaml:"caches"`
	Mounts   []Mount           `yaml:"mounts"`
	Env      map[string]string `yaml:"env"`
	Workdir  string            `yaml:"workdir"`
	Commands []string          `yaml:"commands"`
}

type Stiletto struct {
	Jobs []Job `yaml:"jobs"`
}
