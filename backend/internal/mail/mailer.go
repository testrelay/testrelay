package mail

import (
	"bytes"
	"fmt"
	"html/template"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/testrelay/testrelay/backend/internal/core"
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
		"reviewer-invite": {
			plain: reviewInvitePlain,
			html:  reviewInviteHTML,
		},
		"candidate-invite": {
			plain: candidateInvitePlain,
			html:  candidateInviteHtml,
		},
	}
)

type t struct {
	plain string
	html  string
}

func buildTemplates(name string, data interface{}) (*bytes.Buffer, *bytes.Buffer, error) {
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

type SMTPMailer struct {
	client *mail.SMTPClient
	Domain string
}

func NewSMTPMailer(config core.SMTPConfig) (SMTPMailer, error) {
	server := mail.NewSMTPClient()

	server.Host = config.Host
	server.Port = config.Port
	server.Username = config.Username
	server.Password = config.Password

	smtpClient, err := server.Connect()
	if err != nil {
		return SMTPMailer{}, fmt.Errorf("could not connect to smtp server %w", err)
	}

	return SMTPMailer{
		client: smtpClient,
		Domain: config.SendingDomain,
	}, nil
}

func (s SMTPMailer) Send(config core.MailConfig, data interface{}) error {
	_, html, err := buildTemplates(config.TemplateName, data)
	if err != nil {
		return fmt.Errorf("could not build templates for test %s %w", config.TemplateName, err)
	}

	from := "info@" + s.Domain
	if config.From != "" {
		from = config.From + s.Domain
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(config.To).
		SetSubject(config.Subject).
		SetBody(mail.TextHTML, html.String())

	return email.Send(s.client)
}
