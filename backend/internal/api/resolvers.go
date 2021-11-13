package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
	hGraph "github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/httputil"
)

// TestRepositoryResolver implements a graphql resolver using a vcs app installation to
// find test repositories for a given organisation.
type TestRepositoryResolver struct {
	HasuraURL string
	Collector core.RepoCollector
}

// ResolveRepos returns a list of test repositories for the provided business_id in the graphql params.
// It expects that a vcs app has been installed on the business and fetches the installation_id from
// storage. ResolveRepos errors if no valid installation can be found or if fetching repositories fails.
func (r *TestRepositoryResolver) ResolveRepos(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["business_id"].(int)
	if !ok {
		return []core.Repo{}, nil
	}

	var q struct {
		BusinessByPK struct {
			GithubInstallationID hGraph.String `graphql:"github_installation_id"`
		} `graphql:"businesses_by_pk(id: $id)"`
	}

	client := hGraph.NewClient(r.HasuraURL,
		&http.Client{
			Transport: &httputil.BearerTransport{Token: fmt.Sprintf("%s", p.Context.Value("token"))},
		},
	)

	err := client.Query(context.Background(), &q, map[string]interface{}{
		"id": hGraph.Int(id),
	})
	if err != nil {
		log.Printf("failed to query hasura with id %d err %s\n", id, err)
		return []core.Repo{}, nil
	}

	if q.BusinessByPK.GithubInstallationID == "" {
		log.Printf("returned nil github installation for business")
		return []core.Repo{}, nil
	}

	installationID := q.BusinessByPK.GithubInstallationID
	in, _ := strconv.ParseInt(string(installationID), 10, 64)

	return r.Collector.CollectRepos(in)
}
