package internal

type CandidateEmailData struct {
	Sender       string
	EmailLink    string
	BusinessName string
	Assignment   Assignment
}

type EmailData struct {
	Sender        string
	EmailLink     string
	BusinessName  string
	Email         string
	CandidateName string
}

type MailConfig struct {
	TemplateName string
	Subject      string
	From         string
	To           string
}
