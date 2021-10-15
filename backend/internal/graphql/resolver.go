package graphql

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/graphql-go/graphql"
	hGraph "github.com/hasura/go-graphql-client"
)

type Repo struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

type RepoResolver struct {
	HasuraClient *hGraph.Client
}

func (r *RepoResolver) ResolveRepos(p graphql.ResolveParams) (interface{}, error) {
	var l []Repo
	id, ok := p.Args["business_id"].(int)
	if !ok {
		return l, nil
	}

	var q struct {
		BusinessByPK struct {
			GithubInstallationID hGraph.String `graphql:"github_installation_id"`
		} `graphql:"businesses_by_pk(id: $id)"`
	}

	err := r.HasuraClient.Query(context.Background(), &q, map[string]interface{}{
		"id": hGraph.Int(id),
	})
	if err != nil {
		log.Printf("failed to query hasura with id %d err %s\n", id, err)
		return l, nil
	}

	if q.BusinessByPK.GithubInstallationID == "" {
		log.Printf("returned nil github installation for business")
		return l, nil
	}

	installationID := q.BusinessByPK.GithubInstallationID
	in, _ := strconv.ParseInt(string(installationID), 10, 64)
	pkey := os.Getenv("GITHUB_PRIVATE_KEY")
	pkey = strings.ReplaceAll(pkey, `\n`, "\n")

	itr, err := ghinstallation.New(http.DefaultTransport, 131386, in, []byte(pkey))
	if err != nil {
		log.Printf("failed to init a new github installation %s\n", err)
		return l, nil
	}

	appClient := github.NewClient(&http.Client{Transport: itr})
	repos, _, err := appClient.Apps.ListRepos(context.Background(), nil)
	if err != nil {
		log.Printf("failed to list repose %s\n", err)
		return l, nil
	}

	var qrepos = make([]Repo, len(repos))
	for i, repo := range repos {
		qrepos[i] = Repo{
			ID:       *repo.ID,
			FullName: *repo.FullName,
		}
	}

	return qrepos, nil
}

