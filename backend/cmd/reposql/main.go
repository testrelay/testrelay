package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
	hGraph "github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal/auth"
	intGraphql "github.com/testrelay/testrelay/backend/internal/graphql"
	http2 "github.com/testrelay/testrelay/backend/internal/http"
)

var (
	QueryNameNotProvided = errors.New("no query was provided in the HTTP body")
	verifier             = auth.FirebaseVerifier{
		ProjectID: "testrelay-323914",
	}
	schema   graphql.Schema
	repoType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Repo",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"full_name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	repoResolver = &intGraphql.RepoResolver{}
)

func init() {
	fields := graphql.Fields{
		"repos": &graphql.Field{
			Type:        graphql.NewList(repoType),
			Description: "Get business repos",
			Args: graphql.FieldConfigArgument{
				"business_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: repoResolver.ResolveRepos,
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	var err error
	schema, err = graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}

type params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func doQuery(p params) (events.APIGatewayProxyResponse, error) {
	response := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.OperationName,
	})

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Could not decode response body %s\n", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseJSON),
		StatusCode: http.StatusOK,
	}, nil
}

func Handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no query is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, QueryNameNotProvided
	}

	var p params
	if err := json.Unmarshal([]byte(request.Body), &p); err != nil {
		log.Printf("Could not decode body %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       string("Unauthorized"),
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	log.Printf("request received %+v\n", p)

	// if introspection let return the results now
	if p.OperationName == "IntrospectionQuery" {
		return doQuery(p)
	}

	authH := request.Headers["Authorization"]
	splitToken := strings.Split(authH, "Bearer ")
	if len(splitToken) < 2 {
		return events.APIGatewayProxyResponse{
			Body:       string("Unauthorized"),
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	jwtToken := splitToken[1]
	err := verifier.Parse(jwtToken)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       string("Unauthorized"),
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	repoResolver.HasuraClient = hGraph.NewClient(
		"https://delicate-gator-74.hasura.app/v1/graphql",
		&http.Client{
			Transport: &http2.BearerTransport{Token: jwtToken},
		},
	)

	return doQuery(p)
}

func main() {
	lambda.Start(Handler)
}
