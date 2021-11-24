package api

//go:generate mockgen -destination mocks/handler.go -package mocks . Verifier
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"

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

// Resolver is an interface type that defines a graphl api resolver.
// It must return a set of graphql fields with fieldtypes and query resolvers initialised.
type Resolver interface {
	// Fields returns first the graphql query fields, second the graphql mutation fields.
	// Either can be nil.
	Fields() (graphql.Fields, graphql.Fields)
}

// newSchema initializes a new graphql schema using the given resolvers to build schema fields.
// Note fields with the same name/key will be overwritten with the last entry.
func newSchema(resolvers ...Resolver) (graphql.Schema, error) {
	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)
	for _, r := range resolvers {
		qs, mus := r.Fields()

		if qs != nil {
			for k, field := range qs {
				queries[k] = field
			}
		}

		if mus != nil {
			for k, field := range mus {
				mutations[k] = field
			}
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: queries}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}),
	})
	if err != nil {
		return graphql.Schema{}, fmt.Errorf("could not generate new schema %w", err)
	}

	return schema, nil
}

type Verifier interface {
	Parse(token string) error
}

// GraphQLQueryHandler is a struct that implements the base http.Handler interface.
// It deals with inbound graphql requests, including introspection. Handing off
// queries to local resolvers.
type GraphQLQueryHandler struct {
	hasuraURL string
	schema    graphql.Schema
	verifier  Verifier
}

// NewGraphQLQueryHandler returns a GraphQLQueryHandler with initialized base schema.
// It returns an error if there is an issue initializing the schema for the handler.
// It accepts any number of graphql Resolvers. These Resolvers must define unique query/mutation keys
// Otherwise they will be overwritten when the schema is initialized.
func NewGraphQLQueryHandler(hasuraURL string, verifier Verifier, resolvers ...Resolver) (*GraphQLQueryHandler, error) {
	schema, err := newSchema(resolvers...)
	if err != nil {
		return nil, err
	}

	return &GraphQLQueryHandler{
		hasuraURL: hasuraURL,
		schema:    schema,
		verifier:  verifier,
	}, nil
}

// ServeHTTP implements the http.Handler interface and deals with inbound graphql requests.
// ServeHTTP expects that graphql queries, outside IntrospectionQueries, contain an Authorization header.
// The authorization header must pass jwt validation checks in order to pass to a query Resolver.
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
