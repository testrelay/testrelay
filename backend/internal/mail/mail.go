package mail

import (
	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/core"
)

type Mailer interface {
	SendReviewerInvite(data EmailData) error
	SendCandidateInviteEmail(data CandidateEmailData) error
	Send(config core.MailConfig, data interface{}) error
	SendEnd(status string, data internal.FullAssignment) error
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
