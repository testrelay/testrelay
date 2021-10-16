package mail

import (
	"github.com/testrelay/testrelay/backend/internal"
)

type Mailer interface {
	SendReviewerInvite(data EmailData) error
	SendCandidateInviteEmail(data CandidateEmailData) error
	Send(config Config, data internal.FullAssignment) error
	SendEnd(status string, data internal.FullAssignment) error
}

type Config struct {
	TemplateName string
	Subject      string
	From         string
	To           string
}

type CandidateEmailData struct {
	Sender       string
	EmailLink    string
	BusinessName string
	Assignment   internal.Assignment
}

type EmailData struct {
	Sender        string
	EmailLink     string
	BusinessName  string
	Email         string
	CandidateName string
}
