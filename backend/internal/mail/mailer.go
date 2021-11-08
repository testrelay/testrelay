package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/testrelay/testrelay/backend/internal/core"
)

var (
	//go:embed templates/*
	templates embed.FS
)

// SMTPMailer implements a mailer interface, sending mail through a
// smtp server. SMTPMailer expects Domain to be configured with the
// correct sending from domain. i.e. @mydomain.io.
type SMTPMailer struct {
	server *mail.SMTPServer
	Domain string
}

// NewSMTPMailer returns a mailer instance with the provided config.
func NewSMTPMailer(config core.SMTPConfig) SMTPMailer {
	server := mail.NewSMTPClient()

	server.Host = config.Host
	server.Port = config.Port
	server.Username = config.Username
	server.Password = config.Password

	return SMTPMailer{
		server: server,
		Domain: config.SendingDomain,
	}
}

// Send implements the mailer.Send interface, connecting to a smtp server to
// send a message with the provided data. It first finds a template with the provided
// config and then passes it the data interface.
func (s SMTPMailer) Send(config core.MailConfig, data interface{}) error {
	html, err := buildTemplate(config.TemplateName, data)
	if err != nil {
		return fmt.Errorf("could not build templates for test %s %w", config.TemplateName, err)
	}

	from := "info" + s.Domain
	if config.From != "" {
		from = config.From + s.Domain
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(config.To).
		SetSubject(config.Subject).
		SetBody(mail.TextHTML, html)

	conn, err := s.server.Connect()
	if err != nil {
		return fmt.Errorf("failed to get smtp server connection %w", err)
	}

	return email.Send(conn)
}

func buildTemplate(name string, data interface{}) (string, error) {
	t, err := template.ParseFS(templates, "templates/layout.html.tmpl", "templates/"+name+".html.tmpl")
	if err != nil {
		return "", fmt.Errorf("error opening email template %s %w", name, err)
	}

	html := bytes.NewBuffer([]byte{})
	err = t.ExecuteTemplate(html, "layout", data)
	if err != nil {
		return "", fmt.Errorf("could not execute html template %w", err)
	}

	return html.String(), nil
}
