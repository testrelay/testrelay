import { gql } from '@apollo/client';

const INSERT_ASSIGNMENT = gql`
    mutation InsertAssignment(
      $time_limit: Int!, 
      $choose_until: date!, 
      $name: String!, 
      $email: String!, 
      $test_id: Int!, 
      $reviewers: [assignment_users_insert_input!]!
    ) {
        insert_assignments_one(
          object: {
            candidate_email: $email, 
            candidate_name: $name, 
            test_id: $test_id, 
            choose_until: $choose_until, 
            time_limit: $time_limit,
            reviewers: { data: $reviewers }
          }
        ) {
            id
            candidate_email
        }
    }
`

const GET_ASSIGNMENTS = gql`
    query GetAssignments($offset: Int = 0, $limit: Int = 10) {
         assignments_aggregate {
            aggregate {
                count(distinct: false)
            }
        }
        assignments(order_by: {created_at: desc}, limit: $limit, offset: $offset) {
            id
            candidate_email
            candidate_name
            status
            choose_until
            time_limit
            created_at
            test {
                id
                name
            }
        }
    }
`

const DELETE_REVIEWER = gql`
mutation DeleteReviewer($user_id: Int!) {
  delete_assignment_users(where: {user_id: {_eq: $user_id}}) {
    affected_rows
  }
}
`

const INSERT_REVIEWER = gql`
mutation InsertReviewer($assignment_id: Int!, $user_id: Int!) {
  insert_assignment_users_one(object: {user_id: $user_id, assignment_id: $assignment_id}) {
    id
  }
}
`

const GET_ASSIGNED = gql`
query GetAssigned($uid: String!) {
  assignment_users(where: {user: {auth_id: {_eq: $uid}}}) {
    assignment {
      test {
        name
      }
      id
      candidate_name
      status
      test_day_chosen
      test_time_chosen
      test_timezone_chosen
      time_limit
      github_repo_url
    }
  }
}
`

const GET_ASSIGNMENT = gql`
query GetAssignment($id: Int!) {
  assignments_by_pk(id: $id) {
    candidate_email
    candidate_id
    candidate_name
    choose_until
    created_at
    github_repo_url
    id
    recruiter_id
    status
    test_day_chosen
    test_id
    test_time_chosen
    test_timezone_chosen
    time_limit
    updated_at
    reviewers {
      user {
        email
        github_username
        id
      }
    }
    assignment_events(order_by: {created_at: asc}) {
      id
      event_type
      meta
      created_at
      user {
        id
        email
      }
    }
    test {
      id
      name
      test_languages {
        language {
          name
          id
        }
      }
      github_repo
      zip
    }
  }
}
`
export { INSERT_ASSIGNMENT, GET_ASSIGNMENTS, GET_ASSIGNED, GET_ASSIGNMENT, INSERT_REVIEWER, DELETE_REVIEWER };