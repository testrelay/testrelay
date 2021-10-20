package core

type MailConfig struct {
	TemplateName string
	Subject      string
	From         string
	To           string
}

type SMTPConfig struct {
	SendingDomain string

	Host     string
	Port     int
	Username string
	Password string
}

type Mailer interface {
	Send(config MailConfig, data interface{}) error
}
