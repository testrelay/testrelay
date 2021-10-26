import { gql } from '@apollo/client';

const GET_REPOS = gql`
    query GetRepos($id: Int!) {
      repos(business_id: $id) {
        full_name
        id
      }
    }
`

const GET_LANGUAGES = gql`
  query GetLanguages {
    languages {
      id
      name
    }
  }
`

const INSERT_LANGUAGE = gql`
  mutation InsertLanguage($name: String!) {
    insert_languages_one(object: {name: $name}) {
      id
      name
    }
  }
`

const INSERT_REPO = gql`
  mutation InsertTest($test_window: Int!, $name: String!, $github_repo: String!, $time_limit: Int!, $zip: String, $languages: [test_languages_insert_input!] = {}, $business_id: Int!) {
    insert_tests_one(object: {github_repo: $github_repo, name: $name, test_window: $test_window, time_limit: $time_limit, zip: $zip, test_languages: {data: $languages}, business_id: $business_id}) {
      id
    }
  }
`

const GET_TESTS = gql`
    query GetTests($offset: Int = 0, $limit: Int = 10) {
        tests(order_by: {created_at: desc}, limit: $limit, offset: $offset) {
            id
            name
            zip
            github_repo
            test_window
            time_limit
            created_at
        }
        tests_aggregate {
            aggregate {
              count(columns: id)
            }
        }
    }
`

const GET_TEST = gql`
query GET_TEST($id: Int!) {
  tests_by_pk(id: $id) {
    assignments {
      id
      candidate_email
      candidate_name
      status
    }
    test_languages {
      language {
        name
      }
    }
    test_window
    time_limit
    zip
    github_repo
    name
  }
}
`

export { GET_REPOS, INSERT_REPO, GET_TESTS, GET_LANGUAGES, INSERT_LANGUAGE, GET_TEST };
