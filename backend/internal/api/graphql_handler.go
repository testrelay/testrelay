package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/testrelay/testrelay/backend/internal/auth"
)

var (
	QueryNameNotProvided = errors.New("no query was provided in the HTTP body")
	verifier             = auth.FirebaseVerifier{
		ProjectID: "testrelay-323914",
	}
)

type QueryRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func NewSchema(resolver RepoResolver) (graphql.Schema, error) {
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

type GraphQLQueryHandler struct {
	hasuraURL string
	schema    graphql.Schema
}

func NewGraphQLQueryHandler(hasuraURL string, resolver RepoResolver) (*GraphQLQueryHandler, error) {
	schema, err := NewSchema(resolver)
	if err != nil {
		return nil, err
	}

	return &GraphQLQueryHandler{
		hasuraURL: hasuraURL,
		schema:    schema,
	}, nil
}

func (h *GraphQLQueryHandler) Query(w http.ResponseWriter, r *http.Request) {
	var qr QueryRequest

	err := json.NewDecoder(r.Body).Decode(&qr)
	if err != nil {
		log.Printf("Could not decode body %s\n", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": ["bad request"] }`))
	}

	// if introspection let return the results now
	if qr.OperationName == "IntrospectionQuery" {
		h.doQuery(context.Background(), qr, w)
		return
	}

	authH := r.Header.Get("Authorization")
	splitToken := strings.Split(authH, "Bearer ")
	if len(splitToken) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors": ["unauthorized"] }`))
		return
	}

	jwtToken := splitToken[1]
	err = verifier.Parse(jwtToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors": ["unauthorized"] }`))
		return
	}

	h.doQuery(context.WithValue(context.Background(), "token", jwtToken), qr, w)
}

func (h *GraphQLQueryHandler) doQuery(ctx context.Context, qr QueryRequest, w http.ResponseWriter) {
	response := graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  qr.Query,
		VariableValues: qr.Variables,
		OperationName:  qr.OperationName,
		Context:        ctx,
	})

	json.NewEncoder(w).Encode(response)
}
