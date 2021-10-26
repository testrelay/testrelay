import { gql } from "@apollo/client";


const GET_USER = gql`
query GetUser($limit: Int = 1) {
    users(limit: $limit) {
      id
      email
      github_username
      github_access_token
      auth_id
    }
}`

export { GET_USER };