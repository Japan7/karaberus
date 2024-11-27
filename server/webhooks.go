package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Webhook struct {
	Type string
	URL  string
}

type WebhookTemplateContext struct {
	Kara        KaraInfoDB
	Server      string
	Title       string
	Description string
	Resource    string
}

func parseWebhooksConfig() []Webhook {
	var webhooks []Webhook
	for _, kv := range CONFIG.Webhooks {
		typ, url, found := strings.Cut(kv, "=")
		if !found {
			getLogger().Printf("invalid webhook value: %s", kv)
			continue
		}
		webhooks = append(webhooks, Webhook{typ, url})
	}
	return webhooks
}

func PostWebhooks(kara KaraInfoDB) {
	title := kara.FriendlyName()
	desc, err := karaDescription(kara)
	if err != nil {
		getLogger().Printf("error generating description for webhooks: %s", err)
		return
	}

	tmplCtx := WebhookTemplateContext{
		Kara:        kara,
		Server:      CONFIG.Listen.Addr(),
		Title:       title,
		Description: desc,
		Resource:    fmt.Sprintf("%s/karaoke/browse/%d", CONFIG.Listen.Addr(), kara.ID),
	}

	for _, webhook := range parseWebhooksConfig() {
		var err error
		switch webhook.Type {
		case "json":
			err = postJsonWebhook(webhook.URL, tmplCtx)
		case "discord":
			err = postDiscordWebhook(webhook.URL, tmplCtx)
		default:
			err = fmt.Errorf("unknown webhook type %s", webhook.Type)
		}
		if err != nil {
			getLogger().Printf("error during %s webhook: %s", webhook, err)
		}
	}
}

func postJsonWebhook(url string, tmplCtx WebhookTemplateContext) error {
	b, err := json.Marshal(tmplCtx)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer Closer(resp.Body)
	return nil
}

type DiscordEmbedAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type DiscordEmbed struct {
	Author      DiscordEmbedAuthor `json:"author,omitempty"`
	Title       string             `json:"title"`
	URL         string             `json:"url"`
	Description string             `json:"description"`
	Color       uint               `json:"color"`
}

type DiscordWebhook struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func postDiscordWebhook(url string, tmplCtx WebhookTemplateContext) error {
	webhook_data := DiscordWebhook{
		Embeds: []DiscordEmbed{DiscordEmbed{
			Author: DiscordEmbedAuthor{
				Name:    "New Karaoke!",
				IconURL: fmt.Sprintf("%s/vite.svg", tmplCtx.Server),
			},
			Title:       tmplCtx.Title,
			URL:         tmplCtx.Resource,
			Description: tmplCtx.Description,
			Color:       10053324,
		}},
	}

	body, err := json.Marshal(webhook_data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer Closer(resp.Body)
	return nil
}
