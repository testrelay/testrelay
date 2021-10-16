package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal/core/assignmentuser"
	"github.com/testrelay/testrelay/backend/internal/httputil"
)

type ReviewerHandler struct {
	Logger   *zap.SugaredLogger
	Assigner assignmentuser.Assigner
}

func (rh ReviewerHandler) EventsHandler(w http.ResponseWriter, r *http.Request) {
	var data HasuraEvent
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		body, _ := ioutil.ReadAll(r.Body)
		rh.Logger.Error(
			"could not reviewer events data",
			"error", err,
			"body", body,
		)

		httputil.BadRequest(w)
		return
	}

	if data.Event.Op == "INSERT" {
		var body assignmentuser.RawReviewer
		if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
			rh.Logger.Error(
				"could not reviewer events body",
				"error", err,
				"data", data.Event.Data.New,
			)

			httputil.BadRequest(w)
		}

		err := rh.Assigner.Assign(body)
		if err != nil {
			rh.Logger.Error(
				"could not assign reviewer",
				"assignment_user", body,
				"error", err,
			)

			httputil.BadRequest(w)
		}
	}
}
