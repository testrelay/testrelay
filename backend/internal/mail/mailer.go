package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/testrelay/testrelay/backend/internal"
)

var (
	reviewInvitePlain = `You've been invited to review %s's technical assignment
Check out all your assigned reviews here: %s
`
	reviewInviteHTML = `
<p>You've been invited to review %s's technical assignment<p>
<p>Check out all your assigned reviews <a href="%s">Click here</a></p>
`

	candidateInvitePlain = `Hello %s,
%s has invited you to take a technical test.
Click here %s to schedule your assignment. The test is %s in length and you have until %s to take the test.
When choosing your preferred date/time to sit the test, you'll be prompted to sign in with github. This is so that we can invite you to a private repo where you'll take the test.
Good luck and feel free to reply to this email if you have any technical problems scheduling your technical test.
- the TestRelay candidate team
`
	candidateInviteHtml = `<html><body><h1>Hello %s,</h1>
<p>%s has invited you to take a technical test.<p>
<p><a href="%s">Click here</a> to schedule your assignment. The test is <b>%s</b> in length and you have until <b>%s</b> to take the test.</p>
<p>When choosing your preferred date/time to sit the test, you'll be prompted to sign in with github. This is so that we can invite you to a private repo where you'll take the test.</p>
<p>Good luck and feel free to reply to this email if you have any technical problems scheduling your technical test.</p>
<p><em>- the TestRelay candidate team</em></p></body></html>
`

	candidateStepInvitePlain = `Hello {{ .CandidateName }}, 
your test is due to start in 5 minutes. Your test instructions will be uploaded here: {{.GithubRepoURL}} `
	candidateStepInviteHtml = `<html>
  <body>
    <h1>Hello {{ .CandidateName }},</h1>
    <p>Your test is due to start in 5 minutes. Your test instructions will be uploaded here <a href="{{.GithubRepoURL}}">{{.GithubRepoURL}}</a></p>
  </body>
</html>
`
	endEmailPlain = `Hello {{ .CandidateName }}, 
Your test is due is about to end. Finish up and commit your final changes.`
	endEmailHtml = `<html>
  <body>
    <h1>Hello {{ .CandidateName }},</h1>
    <p>Your test is due is about to end. Finish up and commit your final changes.</p>
  </body>
</html>`

	submittedPlain = `Hello {{ .CandidateName }},
Thanks for submitting your assignment. {{.Test.Business.Name}} will now review your code and get back to you with feedback.`
	submittedHTML = `<html>
<body>
<h1>Hello {{ .CandidateName }},</h1>
<p>Thanks for submitting your assignment. {{.Test.Business.Name}} will now review your code and get back to you with feedback.</p>
</body>
</html>`

	missedPlain = `Hello {{ .CandidateName }},
unfortunately you missed the deadline to submit your assignment. In many cases this is an automatic fail, but you should reach out to {{.Test.Business.Name}} to make sure.`
	missedHtml = `<html>
<body>
<h1>Hello {{ .CandidateName }},</h1>
<p>unfortunately you missed the deadline to submit your assignment. In many cases this is an automatic fail, but you should reach out to {{.Test.Business.Name}} to make sure.</p>
</body>
</html>`

	submittedPlainR = `Candidate {{ .CandidateName }},
Has submitted their assignment. You can check it out here: {{.GithubRepoURL}}.`
	submittedHTMLR = `<html>
<body>
<h1>Good news,</h1>
<p>Candidate {{ .CandidateName }},Has submitted their assignment. You can check it out here: <a href="{{.GithubRepoURL}}">{{.GithubRepoURL}}</a>.</p>
</body>
</html>`

	missedPlainR = `{{ .CandidateName }} missed the deadline to submit their assignment. Try reaching out to them to understand why they weren't able to complete the assignment.`
	missedHtmlR  = `<html>
<p>{{ .CandidateName }} missed the deadline to submit their assignment. Try reaching out to them to understand why they weren't able to complete the assignment.</p>
</body>
</html>`

	templates = map[string]t{
		"warning": {
			plain: candidateStepInvitePlain,
			html:  candidateStepInviteHtml,
		},
		"end": {
			plain: endEmailPlain,
			html:  endEmailHtml,
		},
		"submitted": {
			plain: submittedPlain,
			html:  submittedHTML,
		},
		"submitted-recruiter": {
			plain: submittedPlainR,
			html:  submittedHTMLR,
		},
		"missed": {
			plain: missedPlain,
			html:  missedHtml,
		},
		"missed-recruiter": {
			plain: missedPlainR,
			html:  missedHtmlR,
		},
	}
)

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

type t struct {
	plain string
	html  string
}

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

type MailgunMailer struct {
	MG *mailgun.MailgunImpl
}

func (m *MailgunMailer) SendEnd(status string, data internal.FullAssignment) error {
	subject := "Thanks for submitting your test for " + data.Test.Business.Name
	if status != "submitted" {
		subject = "You missed the deadline for submitting your test"
	}

	err := m.Send(Config{
		TemplateName: status,
		Subject:      subject,
		From:         "candidates@testrelay.io",
		To:           data.CandidateEmail,
	}, data)
	if err != nil {
		return fmt.Errorf("could not send email to candidate %w", err)
	}

	subject = data.CandidateName + " has submitted their assignment"
	if status != "submitted" {
		subject = data.CandidateName + " missed the deadline to submit their assignment"
	}

	err = m.Send(Config{
		TemplateName: status + "-recruiter",
		Subject:      subject,
		From:         "candidates@testrelay.io",
		To:           data.Recruiter.Email,
	}, data)
	if err != nil {
		return fmt.Errorf("could not send email to recruiter %w", err)
	}

	return err
}

func (m *MailgunMailer) Send(config Config, data internal.FullAssignment) error {
	plain, html, err := m.buildTemplates(config.TemplateName, data)
	if err != nil {
		return fmt.Errorf("could not build templates for test %s %w", config.TemplateName, err)
	}

	message := m.MG.NewMessage(
		config.From,
		config.Subject,
		plain.String(),
		config.To,
	)
	message.SetHtml(html.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = m.MG.Send(ctx, message)
	return err
}

func (m *MailgunMailer) buildTemplates(name string, data interface{}) (*bytes.Buffer, *bytes.Buffer, error) {
	t := template.New("plain")
	t, _ = t.Parse(templates[name].plain)

	ht := template.New("html")
	ht, _ = ht.Parse(templates[name].html)

	plain := bytes.NewBuffer([]byte{})
	err := t.Execute(plain, data)
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute plain template %w", err)
	}

	html := bytes.NewBuffer([]byte{})
	err = ht.Execute(html, data)
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute html template %w", err)
	}

	return plain, html, nil
}

func (m *MailgunMailer) SendReviewerInvite(data EmailData) error {
	message := m.MG.NewMessage(
		data.Sender,
		"You've been invited you to review "+data.CandidateName+"'s technical assignment",
		m.invitePlain(data),
		data.Email,
	)
	message.SetHtml(m.inviteHtml(data))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := m.MG.Send(ctx, message)
	return err
}

func (m *MailgunMailer) inviteHtml(data EmailData) string {
	return fmt.Sprintf(
		reviewInviteHTML,
		data.CandidateName,
		fmt.Sprintf("%s/reviews", os.Getenv("APP_URL")),
	)
}

func (m *MailgunMailer) invitePlain(data EmailData) string {
	return fmt.Sprintf(
		reviewInvitePlain,
		data.CandidateName,
		fmt.Sprintf("%s/reviews", os.Getenv("APP_URL")),
	)
}

func (m *MailgunMailer) SendCandidateInviteEmail(data CandidateEmailData) error {
	message := m.MG.NewMessage(
		data.Sender,
		data.BusinessName+" has invited you to a technical test",
		m.candidateInvitePlain(data),
		data.Assignment.CandidateEmail,
	)
	message.SetHtml(m.candidateInviteHtml(data))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := m.MG.Send(ctx, message)
	return err
}

func (m *MailgunMailer) candidateInviteHtml(data CandidateEmailData) string {
	return fmt.Sprintf(
		candidateInviteHtml,
		data.Assignment.CandidateName,
		data.BusinessName,
		data.EmailLink,
		testTime(data.Assignment.TimeLimit),
		chooseUntil(data.Assignment.ChooseUntil),
	)
}

func (m *MailgunMailer) candidateInvitePlain(data CandidateEmailData) string {
	return fmt.Sprintf(
		candidateInvitePlain,
		data.Assignment.CandidateName,
		data.BusinessName,
		data.EmailLink,
		testTime(data.Assignment.TimeLimit),
		chooseUntil(data.Assignment.ChooseUntil),
	)
}

func chooseUntil(until string) string {
	t, _ := time.Parse("2006-01-02", until)
	return t.Format("Mon, 02 Jan 2006")
}

func testTime(limit int) string {
	hours := limit / (60 * 60)

	if hours == 1 {
		return fmt.Sprintf("%d hour", hours)
	}

	if hours <= 24 {
		return fmt.Sprintf("%d hours", hours)
	}

	days := hours / 24
	return fmt.Sprintf("%d days", days)
}
