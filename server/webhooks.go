package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Webhook struct {
	Type string
	URL  string
}

type WebhookTemplateContext struct {
	Config      *KaraberusConfig
	Kara        *KaraInfoDB
	Title       string
	Description string
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

func PostWebhooks(kara *KaraInfoDB) error {
	if kara == nil {
		return nil
	}
	title := kara.FriendlyName()
	desc, err := karaDescription(*kara)
	if err != nil {
		return err
	}
	tmplCtx := &WebhookTemplateContext{&CONFIG, kara, title, desc}
	for _, webhook := range parseWebhooksConfig() {
		var err error
		switch webhook.Type {
		case "json":
			err = postJsonWebhook(webhook.URL, tmplCtx)
		case "discord":
			err = postDiscordWebhook(webhook.URL, tmplCtx)
		default:
			err = fmt.Errorf("unknown webhook type: %s", webhook.Type)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func postJsonWebhook(url string, tmplCtx *WebhookTemplateContext) error {
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

func postDiscordWebhook(url string, tmplCtx *WebhookTemplateContext) error {
	tmpl, err := template.New("discord_template").Funcs(funcs).Parse(discordTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplCtx); err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)
	return nil
}

var funcs = template.FuncMap{
	"jsstr": func(s string) template.JSStr {
		return template.JSStr(strings.ReplaceAll(s, "\n", "\\n"))
	},
	"join": func(sep string, elems []string) string {
		return strings.Join(elems, sep)
	},
}

const discordTemplate = `{
	"embeds": [
		{
			"author": {
				"name": "New Karaoke!",
				"icon_url": "{{ .Config.Listen.BaseURL }}/vite.svg"
			},
			"title": "{{ .Title | jsstr }}",
			"url": "{{ .Config.Listen.BaseURL }}/karaoke/browse/{{ .Kara.ID }}",
			"description":  "{{ .Description | jsstr }}",
			"color": 10053324
		}
	]
}`
