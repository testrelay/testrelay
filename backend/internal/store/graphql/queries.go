package graphql

import (
	"github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal"
)

type UserBQuery struct {
	TestByPk struct {
		Business struct {
			Name graphql.String `graphql:"name"`
		} `graphql:"business"`
	} `graphql:"tests_by_pk(id: $test_id)"`
	Users []struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"users(where: {auth_id: {_eq: $user_id}})"`
}

type UserQuery struct {
	Users []struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"users(where: {auth_id: {_eq: $user_id}})"`
}

type InsertUserMutation struct {
	InsertUsersOne struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"insert_users_one(object: {auth_id: $auth_id, email: $email})"`
}

type BusinessQuery struct {
	TestByPk struct {
		Business struct {
			Name graphql.String `graphql:"name"`
			ID   graphql.Int    `graphql:"id"`
		} `graphql:"business"`
	} `graphql:"tests_by_pk(id: $test_id)"`
}

type UpdateAssignmentMutation struct {
	UpdateAssignmentsByPK struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_assignments_by_pk(pk_columns: {id: $id}, _set: {status: $status, candidate_id: $candidate_id})"`
	InsertAssignmentEvents struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"insert_assignment_events(objects: {event_type: $status, user_id: $user_id, assignment_id: $id})"`
	InsertBusinessUsersOne struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"insert_business_users_one(object: {business_id: $business_id, user_id: $candidate_id, user_type: $user_type},on_conflict: {constraint: business_users_business_id_user_id_user_type_key})"`
}

type InsertAssignmentEvent struct {
	UpdateAssignmentsByPK struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_assignments_by_pk(pk_columns: {id: $id}, _set: {status: $status})"`
	InsertAssignmentEvents struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"insert_assignment_events(objects: {event_type: $status, user_id: $user_id, assignment_id: $id})"`
}

type AssignmentReviewers struct {
	AssignmentUsers struct {
		Reviewers []internal.Reviewer `graphql:"reviewers"`
	} `graphql:"assignments_by_pk(id: $id)"`
}

type Reviewer struct {
	GithubUsername string `graphql:"github_username" json:"github_username"`
}

type assignment_status_enum string

func newStatus(s string) *assignment_status_enum {
	a := assignment_status_enum(s)
	return &a
}
