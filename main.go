package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"strings"
)

const (
	GoodTrue  = "✅"
	BadTrue   = "☑️"
	GoodFalse = "❎"
	BadFalse  = "❌"
)

func Select(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func main() {
	data, err := ioutil.ReadFile("clients.yaml")
	if err != nil {
		panic(err)
	}

	var config clientsConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	// build target key map
	identifierTargetMap := make(map[string]*Target)
	for _, target := range config.Targets {
		for has := range target.Has {
			if _, ok := identifierTargetMap[has]; ok {
				panic(has + " specified twice!")
			}
			identifierTargetMap[has] = target
		}
	}

	identifierClientMap := make(map[string][]*Client)
	for _, client := range config.Clients {
		for _, targetStr := range client.TargetsStr {
			if _, ok := identifierTargetMap[targetStr]; !ok {
				panic("cannot find target " + targetStr)
			}
			identifierClientMap[targetStr] = append(identifierClientMap[targetStr], client)
		}
	}

	for _, target := range config.Targets {
		fmt.Printf("## %s\n\n", target.Display)
		for has, hasDisplay := range target.Has {
			fmt.Printf("### %s\n\n", hasDisplay)

			fmt.Println("| Name | Website | OSS | Free | Paid | Downloads |")
			fmt.Println("|------|---------|-----|------|------|-----------|")

			for _, client := range identifierClientMap[has] {
				var (
					name                      string
					websiteURL, websiteDomain string
					oss, free, paid           string
					downloads                 strings.Builder
				)

				name = client.Name
				if client.Official {
					name += " 🔹"
				}
				free = Select(client.Price.Free, GoodTrue, BadFalse)
				paid = Select(client.Price.Paid, BadTrue, GoodFalse)
				oss = Select(client.OpenSourceURL != "", GoodTrue, BadFalse)

				for hoster, u := range client.Downloads {
					if downloads.Len() > 0 {
						downloads.WriteString(" ")
					}

					dl, ok := config.Downloads[hoster]
					if !ok {
						panic("cannot find download " + hoster)
					}

					img := fmt.Sprintf(`
<a href="%s">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="%s">
		<source media="(prefers-color-scheme: light)" srcset="%s">
		<img src="%s">
	</picture>
</a>`, u, dl.Icons.Light, dl.Icons.Dark, dl.Icons.Dark)

					downloads.WriteString(strings.ReplaceAll(strings.ReplaceAll(img, "\n", ""), "\t", ""))
				}

				if client.OpenSourceURL != "" {
					websiteURL = client.OpenSourceURL
				}
				if client.Website != "" {
					websiteURL = client.Website
				}
				if websiteURL != "" {
					if u, err := url.Parse(websiteURL); err != nil {
						panic(err)
					} else {
						websiteDomain = u.Host
					}
				}

				fmt.Printf("| %s | [%s](%s) | %s | %s | %s | %s |\n",
					name, websiteDomain, websiteURL, oss, free, paid, downloads.String())
			}
			fmt.Println()
		}
	}
}
