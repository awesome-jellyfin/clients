package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	GoodTrue    = "✅"
	BadTrue     = "☑️"
	GoodFalse   = "❎"
	BadFalse    = "❌"
	TypeMusic   = "🎵"
	Official    = "🔹"
	Beta        = "🛠️"
	JellyfinOrg = "https://github.com/jellyfin/"
)

func Select[T any](expr bool, whenTrue, whenFalse T) T {
	if expr {
		return whenTrue
	}
	return whenFalse
}

func Ref[T any](what T) *T {
	return &what
}

func DerefDef[T any](what *T, defaultValue T) T {
	if what != nil {
		return *what
	}
	return defaultValue
}

func Deref[T any](what *T) T {
	if what != nil {
		return *what
	}
	var def T
	return def
}

func main() {
	data, err := os.ReadFile("clients.yaml")
	if err != nil {
		panic(err)
	}
	var config clientsConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	// target -> list of clients
	identifierClientMap := make(map[string][]*Client)
	for _, client := range config.Clients {
		for _, targetStr := range client.Targets {
			targetStr = strings.TrimSpace(strings.ToLower(targetStr))
			identifierClientMap[targetStr] = append(identifierClientMap[targetStr], client)
		}
	}
	for _, target := range config.Targets {
		fmt.Printf("## %s\n\n", target.Display)
		for has, hasDisplay := range target.Has {
			fmt.Printf("### %s\n\n", hasDisplay)

			fmt.Println("| Name | OSS | Free | Paid | Downloads |")
			fmt.Println("|------|-----|------|------|-----------|")

			for _, client := range identifierClientMap[strings.TrimSpace(strings.ToLower(has))] {
				var (
					name            = client.Name
					websiteURL      string
					oss, free, paid string
					downloads       strings.Builder
				)

				// if open-source and in the 'jellyfin' GitHub org, mark as official by default
				if client.Official == nil && strings.HasPrefix(client.OpenSourceURL, JellyfinOrg) {
					client.Official = Ref(true)
				}

				// make free the default price if open-source
				if client.Price.Free == nil && client.OpenSourceURL != "" {
					client.Price.Free = Ref(true)
				}

				/// Pledges
				// append official badge
				if Deref(client.Official) {
					name += " " + Official
				}
				// append beta badge
				if Deref(client.Beta) {
					name += " " + Beta
				}
				// append type badges
				for _, t := range client.Types {
					if t == "Music" {
						name += " " + TypeMusic
					}
				}

				free = Select(DerefDef(client.Price.Free, false), GoodTrue, BadFalse)
				paid = Select(DerefDef(client.Price.Paid, false), BadTrue, GoodFalse)
				oss = Select(client.OpenSourceURL != "", GoodTrue, BadFalse)

				if client.OpenSourceURL != "" {
					websiteURL = client.OpenSourceURL
				}
				if client.Website != "" {
					websiteURL = client.Website
				}

				for _, hoster := range client.Downloads {
					if downloads.Len() > 0 {
						downloads.WriteString(" ")
					}

					// Simple Text Download
					if hoster.Text != "" {
						downloads.WriteString(fmt.Sprintf("[%s](%s)", hoster.Text, hoster.URL))
						continue
					}

					// Icon Download
					var downloadMarkdown string
					if hoster.Icon != "" {
						if icon, ok := config.Icons[hoster.Icon]; ok {
							downloadMarkdown = icon.Markdown(hoster.URL)
						} else {
							panic("cannot find icon " + hoster.Icon + " in config")
						}
					} else if hoster.IconURL != "" {
						downloadMarkdown = (&icon{Single: hoster.IconURL}).Markdown(hoster.URL)
					} else if hoster.Text != "" {
						downloadMarkdown = fmt.Sprintf("[%s](%s)", hoster.Text, hoster.URL)
					} else {
						panic("invalid download. specify either icon, icon-url or text")
					}

					downloads.WriteString(strings.ReplaceAll(
						strings.ReplaceAll(downloadMarkdown, "\n", ""), "\t", ""))
				}

				fmt.Printf("| [%s](%s) | %s | %s | %s | %s |\n",
					name, websiteURL, oss, free, paid, downloads.String())
			}
			fmt.Println()
		}
	}
}
