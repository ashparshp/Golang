package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates

var emailTemplates embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data any) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTemplates, templateToRender)
	if err != nil {
		return fmt.Errorf("error parsing email template %s: %w", templateToRender, err)
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", data); err != nil {
		return fmt.Errorf("error executing email template %s: %w", templateToRender, err)
	}

	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplates, templateToRender)
	if err != nil {
		return fmt.Errorf("error parsing email template %s: %w", templateToRender, err)
	}

	var plainTpl bytes.Buffer
	if err := t.ExecuteTemplate(&plainTpl, "body", data); err != nil {
		return fmt.Errorf("error executing email template %s: %w", templateToRender, err)
	}

	plainMessage := plainTpl.String()

	return nil
}
