package main

type Target struct {
	Key     string            `yaml:"key"`
	Display string            `yaml:"display"`
	Has     map[string]string `yaml:"has"`
}

type Price struct {
	Free bool `yaml:"free"`
	Paid bool `yaml:"paid"`
}

type Client struct {
	Name          string            `yaml:"name"`
	TargetsStr    []string          `yaml:"targets"`
	Official      bool              `yaml:"official"`
	Website       string            `yaml:"website"`
	OpenSourceURL string            `yaml:"oss"`
	Price         *Price            `yaml:"price"`
	Downloads     map[string]string `yaml:"downloads"`
}

type clientsConfig struct {
	Targets   []*Target            `yaml:"targets"`
	Downloads map[string]*Download `yaml:"downloads"`
	Clients   []*Client            `yaml:"clients"`
}

type Icons struct {
	Light string
	Dark  string
}

type Download struct {
	Icons *Icons
}
