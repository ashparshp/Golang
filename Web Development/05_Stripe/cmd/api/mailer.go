package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
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

	app.infoLog.Println(formattedMessage, plainMessage)

	server := mail.NewSMTPClient()
	server.Host = app.config.smtp.host
	server.Port = app.config.smtp.port
	server.Username = app.config.smtp.username
	server.Password = app.config.smtp.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to SMTP server: %w", err)
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	app.infoLog.Printf("Email sent successfully to %s with subject %s", to, subject)

	return nil
}
