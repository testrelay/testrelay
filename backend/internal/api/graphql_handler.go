package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/testrelay/testrelay/backend/internal/auth"
	"github.com/testrelay/testrelay/backend/internal/httputil"
)

const (
	introspectionOp = "IntrospectionQuery"
)

type queryRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// NewSchema initializes a new graphql schema using resolver
// as the base graphql resolver. It returns an error if there
// is an error initializing.
func NewSchema(resolver *TestRepositoryResolver) (graphql.Schema, error) {
	repoType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Repo",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"full_name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	fields := graphql.Fields{
		"repos": &graphql.Field{
			Type:        graphql.NewList(repoType),
			Description: "Get business repos",
			Args: graphql.FieldConfigArgument{
				"business_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: resolver.ResolveRepos,
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return graphql.Schema{}, fmt.Errorf("could not generate new schema %w", err)
	}

	return schema, nil
}

// GraphQLQueryHandler is a struct that implements the base http.Handler interface.
// It deals with inbound graphql requests, including introspection. Handing off
// queries to local resolvers.
type GraphQLQueryHandler struct {
	hasuraURL string
	schema    graphql.Schema
	verifier  auth.FirebaseVerifier
}

// NewGraphQLQueryHandler returns a GraphQLQueryHandler with initialized base schema.
// It returns an error if there is an issue initializing the schema for the handler.
func NewGraphQLQueryHandler(hasuraURL, projectID string, resolver *TestRepositoryResolver) (*GraphQLQueryHandler, error) {
	schema, err := NewSchema(resolver)
	if err != nil {
		return nil, err
	}

	return &GraphQLQueryHandler{
		hasuraURL: hasuraURL,
		schema:    schema,
		verifier: auth.FirebaseVerifier{
			ProjectID: projectID,
		},
	}, nil
}

// ServeHTTP implements the http.Handler interface and deals with inbound graphql requests.
// ServeHTTP expects that graphql queries, outside IntrospectionQueries, contain an Authorization header.
func (h *GraphQLQueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var qr queryRequest

	err := json.NewDecoder(r.Body).Decode(&qr)
	if err != nil {
		log.Printf("could not decode body %s\n", err)

		httputil.BadRequest(w)
		return
	}

	if qr.OperationName == introspectionOp {
		h.doQuery(context.Background(), qr, w)
		return
	}

	authH := r.Header.Get("Authorization")
	splitToken := strings.Split(authH, "Bearer ")
	if len(splitToken) < 2 {
		httputil.Unauthorized(w)
		return
	}

	jwtToken := splitToken[1]
	err = h.verifier.Parse(jwtToken)
	if err != nil {
		httputil.Unauthorized(w)
		return
	}

	h.doQuery(context.WithValue(context.Background(), "token", jwtToken), qr, w)
}

func (h *GraphQLQueryHandler) doQuery(ctx context.Context, qr queryRequest, w http.ResponseWriter) {
	response := graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  qr.Query,
		VariableValues: qr.Variables,
		OperationName:  qr.OperationName,
		Context:        ctx,
	})

	json.NewEncoder(w).Encode(response)
}
