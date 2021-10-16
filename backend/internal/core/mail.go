package core

type MailConfig struct {
	TemplateName string
	Subject      string
	From         string
	To           string
}

type Mailer interface {
	Send(config MailConfig, data interface{}) error
}
