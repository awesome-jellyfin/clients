package main

import "fmt"

// client structs

type Price struct {
	Free *bool `yaml:"free"`
	Paid *bool `yaml:"paid"`
}

type Hoster struct {
	Icon    string `yaml:"icon"`
	IconURL string `yaml:"icon-url"`
	Text    string `yaml:"text"`
	URL     string `yaml:"url"`
}

type Client struct {
	Name          string    `yaml:"name"`
	Targets       []string  `yaml:"targets"`
	Types         []string  `yaml:"types"`
	Official      *bool     `yaml:"official"`
	Beta          *bool     `yaml:"beta"`
	Website       string    `yaml:"website"`
	OpenSourceURL string    `yaml:"oss"`
	Price         Price     `yaml:"price"`
	Downloads     []*Hoster `yaml:"downloads"`
}

// misc

type icon struct {
	Light  string `yaml:"light"`
	Dark   string `yaml:"dark"`
	Single string `yaml:"single"`
	Text   string `yaml:"text"`
}

func (i *icon) Markdown(url string) string {
	if (i.Dark != "") != (i.Light != "") {
		panic("use 'single' if only single icon URL available")
	}
	if i.Dark != "" {
		return fmt.Sprintf(`
<a href="%s">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="%s">
		<source media="(prefers-color-scheme: light)" srcset="%s">
		<img src="%s">
	</picture>
</a>`, url, i.Light, i.Dark, i.Dark)
	}
	if i.Text != "" {
		return fmt.Sprintf("[%s](%s)", i.Text, url)
	}
	return fmt.Sprintf("[![img](%s)](%s)", i.Single, url)
}

type clientsConfig struct {
	Clients []*Client `yaml:"clients"`
	Targets []*struct {
		Key     string            `yaml:"key"`
		Display string            `yaml:"display"`
		Has     map[string]string `yaml:"has"`
	} `yaml:"targets"`
	Icons map[string]*icon `yaml:"icons"`
}
