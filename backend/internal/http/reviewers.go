package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/assignment"
	"github.com/testrelay/testrelay/backend/internal/github"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	"github.com/testrelay/testrelay/backend/internal/mail"
)

type ReviewerHandler struct {
	Logger       *zap.SugaredLogger
	Client       *graphql.HasuraClient
	GithubClient *github.Client
	Mailer       mail.Mailer
}

func (rh ReviewerHandler) EventsHandler(w http.ResponseWriter, r *http.Request) {
	var data assignment.HasuraEvent
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		body, _ := ioutil.ReadAll(r.Body)
		rh.Logger.Error(
			"could not reviewer events data",
			"error", err,
			"body", body,
		)

		BadRequest(w)
		return
	}

	if data.Event.Op == "INSERT" {
		var body internal.AssignmentUser
		if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
			rh.Logger.Error(
				"could not reviewer events body",
				"error", err,
				"data", data.Event.Data.New,
			)

			BadRequest(w)
		}

		au, err := rh.Client.GetAssignmentUser(body.ID)
		if err != nil {
			rh.Logger.Error(
				"could not get assignment user",
				"assignment_user", body.ID,
				"error", err,
			)

			BadRequest(w)
		}

		if au.Assignment.GithubRepoUrl != "" && au.User.GithubUsername != "" {
			err := rh.GithubClient.AddCollaborator(string(au.Assignment.GithubRepoUrl), string(au.User.GithubUsername))
			if err != nil {
				rh.Logger.Error(
					"could not add collaborator",
					"repo_url", au.Assignment.GithubRepoUrl,
					"github_username", au.User.GithubUsername,
					"error", err,
				)

				BadRequest(w)
			}
		}

		data := mail.EmailData{
			Sender:        "info@testrelay.io",
			Email:         string(au.User.Email),
			CandidateName: string(au.Assignment.CandidateName),
		}
		err = rh.Mailer.SendReviewerInvite(data)
		if err != nil {
			rh.Logger.Error(
				"could not send reviewer invite",
				"data", data,
				"error", err,
			)
			BadRequest(w)
		}
	}
}
