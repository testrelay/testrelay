package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
	hGraph "github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal/github"
)

type RepoResolver interface {
	ResolveRepos(p graphql.ResolveParams) (interface{}, error)
}

type RepoCollector interface {
	CollectRepos(installationID int64) ([]github.Repo, error)
}

type GraphResolver struct {
	HasuraURL string
	Collector RepoCollector
}

func (r *GraphResolver) ResolveRepos(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["business_id"].(int)
	if !ok {
		return []github.Repo{}, nil
	}

	var q struct {
		BusinessByPK struct {
			GithubInstallationID hGraph.String `graphql:"github_installation_id"`
		} `graphql:"businesses_by_pk(id: $id)"`
	}

	client := hGraph.NewClient(r.HasuraURL,
		&http.Client{
			Transport: &BearerTransport{Token: fmt.Sprintf("%s", p.Context.Value("token"))},
		},
	)

	err := client.Query(context.Background(), &q, map[string]interface{}{
		"id": hGraph.Int(id),
	})
	if err != nil {
		log.Printf("failed to query hasura with id %d err %s\n", id, err)
		return []github.Repo{}, nil
	}

	if q.BusinessByPK.GithubInstallationID == "" {
		log.Printf("returned nil github installation for business")
		return []github.Repo{}, nil
	}

	installationID := q.BusinessByPK.GithubInstallationID
	in, _ := strconv.ParseInt(string(installationID), 10, 64)

	return r.Collector.CollectRepos(in)
}
