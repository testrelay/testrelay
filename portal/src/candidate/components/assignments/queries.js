import { gql } from '@apollo/client';

const GET_ASSIGNMENT = gql`
query GetAssignment($id: Int!) {
  assignments_by_pk(id: $id) {
    id
    candidate_id
    test_time_chosen
    test_day_chosen
    test_timezone_chosen
    time_limit
    choose_until
    test {
      business {
        name
      }
      test_languages {
          language {
              name
          }
      }
    }
    created_at
  }
}
`

const GET_ASSIGNMENTS = gql`
query GetAssignments {
  assignments {
    github_repo_url
    status
    test_day_chosen
    test_time_chosen
    test_timezone_chosen
    choose_until
    time_limit
    created_at
    id
    test {
      business {
        name
      }
    }
  }
}
`

const CHECK_ASSIGNMENT = gql`
query CheckAssignment($id: Int!, $uuid: uuid!) {
  assignments(where: {_and: {id: {_eq: $id}, invite_code: {_eq: $uuid}, candidate_id: {_is_null: true}}}) {
    id
  }
}
`

const UPDATE_ASSIGNMENT_CANDIDATE = gql`
mutation UpdateAssignment($id: Int!) {
  update_assignments_by_pk(pk_columns: {id: $id}, _set: {status: viewed}) {
    candidate_id
  }
  insert_assignment_events(objects: {event_type: viewed, assignment_id: $id}) {
    affected_rows
  }
}
`

const UPDATE_ASSIGNMENT_CANCELED = gql`
mutation UpdateAssignment($id: Int!) {
  update_assignments_by_pk(pk_columns: {id: $id}, _set: {status: cancelled}) {
    candidate_id
  }
  insert_assignment_events(objects: {event_type: cancelled, assignment_id: $id}) {
    affected_rows
  }
}
`

const UPDATE_ASSIGNMENT_WITH_SCHEDULE = gql`
mutation UpdateAssignmentScheduled($id: Int!, $meta: jsonb!, $test_day_chosen: date!, $test_time_chosen: time!, $test_timezone_chosen: String!) {
  update_assignments_by_pk(pk_columns: {id: $id}, _set: {status: scheduled, test_day_chosen: $test_day_chosen, test_time_chosen: $test_time_chosen, test_timezone_chosen: $test_timezone_chosen}) {
    candidate_id
  }
  insert_assignment_events(objects: {event_type: scheduled, assignment_id: $id, meta: $meta}) {
    affected_rows
  }
}

`

export { GET_ASSIGNMENT, CHECK_ASSIGNMENT, UPDATE_ASSIGNMENT_CANDIDATE, UPDATE_ASSIGNMENT_WITH_SCHEDULE, UPDATE_ASSIGNMENT_CANCELED, GET_ASSIGNMENTS };