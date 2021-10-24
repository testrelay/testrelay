package mail

import (
	"bytes"
	"fmt"
	"html/template"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/testrelay/testrelay/backend/internal/core"
)

var (
	reviewInviteHTML = `
<p>You've been invited to review %s's technical assignment<p>
<p>Check out all your assigned reviews <a href="%s">Click here</a></p>
`

	candidateInviteHtml = `<html><body><h1>Hello {{ .Assignment.CandidateName }},</h1>
<p>{{ .BusinessName }} has invited you to take a technical test.<p>
<p><a href="{{.EmailLink}}">Click here</a> to schedule your assignment. The test is <b>{{ .Assignment.TimeLimitReadable }}</b> in length and you have until <b>{{.Assignment.ChooseReadable}}</b> to take the test.</p>
<p>When choosing your preferred date/time to sit the test, you'll be prompted to sign in with github. This is so that we can invite you to a private repo where you'll take the test.</p>
<p>Good luck and feel free to reply to this email if you have any technical problems scheduling your technical test.</p>
<p><em>- the TestRelay candidate team</em></p></body></html>
`

	candidateStepInviteHtml = `<html>
  <body>
    <h1>Hello {{ .CandidateName }},</h1>
    <p>Your test is due to start in 5 minutes. Your test instructions will be uploaded here <a href="{{.GithubRepoURL}}">{{.GithubRepoURL}}</a></p>
  </body>
</html>
`
	endEmailHtml = `<html>
  <body>
    <h1>Hello {{ .CandidateName }},</h1>
    <p>Your test is due is about to end. Finish up and commit your final changes.</p>
  </body>
</html>`

	submittedHTML = `<html>
<body>
<h1>Hello {{ .CandidateName }},</h1>
<p>Thanks for submitting your assignment. {{.Test.Business.Name}} will now review your code and get back to you with feedback.</p>
</body>
</html>`

	missedHtml = `<html>
<body>
<h1>Hello {{ .CandidateName }},</h1>
<p>unfortunately you missed the deadline to submit your assignment. In many cases this is an automatic fail, but you should reach out to {{.Test.Business.Name}} to make sure.</p>
</body>
</html>`

	submittedHTMLR = `<html>
<body>
<h1>Good news,</h1>
<p>Candidate {{ .CandidateName }},Has submitted their assignment. You can check it out here: <a href="{{.GithubRepoURL}}">{{.GithubRepoURL}}</a>.</p>
</body>
</html>`

	missedHtmlR = `<html>
<p>{{ .CandidateName }} missed the deadline to submit their assignment. Try reaching out to them to understand why they weren't able to complete the assignment.</p>
</body>
</html>`

	templates = map[string]t{
		"warning": {
			html: candidateStepInviteHtml,
		},
		"end": {
			html: endEmailHtml,
		},
		"submitted": {
			html: submittedHTML,
		},
		"submitted-recruiter": {
			html: submittedHTMLR,
		},
		"missed": {
			html: missedHtml,
		},
		"missed-recruiter": {
			html: missedHtmlR,
		},
		"reviewer-invite": {
			html: reviewInviteHTML,
		},
		"candidate-invite": {
			html: candidateInviteHtml,
		},
	}
)

type t struct {
	html string
}

func buildTemplates(name string, data interface{}) (*bytes.Buffer, error) {
	ht := template.New("html")
	ht, _ = ht.Parse(templates[name].html)

	html := bytes.NewBuffer([]byte{})
	err := ht.Execute(html, data)
	if err != nil {
		return nil, fmt.Errorf("could not execute html template %w", err)
	}

	return html, nil
}

type SMTPMailer struct {
	server *mail.SMTPServer
	Domain string
}

func NewSMTPMailer(config core.SMTPConfig) (SMTPMailer, error) {
	server := mail.NewSMTPClient()

	server.Host = config.Host
	server.Port = config.Port
	server.Username = config.Username
	server.Password = config.Password

	return SMTPMailer{
		server: server,
		Domain: config.SendingDomain,
	}, nil
}

func (s SMTPMailer) Send(config core.MailConfig, data interface{}) error {
	html, err := buildTemplates(config.TemplateName, data)
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
		SetBody(mail.TextHTML, html.String())

	conn, err := s.server.Connect()
	if err != nil {
		return fmt.Errorf("failed to get smtp server connection %w", err)
	}

	return email.Send(conn)
}
