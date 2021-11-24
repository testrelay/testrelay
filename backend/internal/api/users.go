package api

//go:generate mockgen -destination mocks/users.go -package mocks . Inviter
import (
	"fmt"

	"github.com/graphql-go/graphql"
	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal/core/user"
)

type Inviter interface {
	Invite(email, redirectLink string, businessID int64) (*user.AuthInfo, error)
}

// UserResolver implements a Resolver interface, declaring methods needed to resolve user queries and mutations.
type UserResolver struct {
	Inviter Inviter
	Logger  *zap.SugaredLogger
}

// Fields returns the queries and mutations defined for user graphql object.
func (u UserResolver) Fields() (graphql.Fields, graphql.Fields) {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	return nil, graphql.Fields{
		"inviteUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"business_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"redirect_link": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: u.InviteUser,
		},
	}
}

// InviteUser parses the graphql request and passes the variables down to the Inviter.
// If there is an error we'll return a friendly user error and handoff the error stack to the logger.
func (u UserResolver) InviteUser(p graphql.ResolveParams) (interface{}, error) {
	email := p.Args["email"].(string)
	var businessId int64
	switch p.Args["business_id"].(type) {
	case int64:
		businessId = p.Args["business_id"].(int64)
	case int:
		businessId = int64(p.Args["business_id"].(int))
	}
	redirectLink := p.Args["redirect_link"].(string)

	a, err := u.Inviter.Invite(email, redirectLink, businessId)
	if err != nil {
		u.Logger.Errorf("could not invite user %s %s", email, err)
		return nil, fmt.Errorf("could not invite user %s to business", email)
	}

	return a, nil
}
