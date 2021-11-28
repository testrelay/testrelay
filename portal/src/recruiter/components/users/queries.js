import {gql} from '@apollo/client';

const GET_USERS = gql`
    query GetUsers($business_id: Int!) {
        users(where: {business_users: {_and: {business_id: {_eq: $business_id}, user_type: {_neq: "candidate"}}}}) {
            email
            github_username
            id
            created_at
            updated_at
            business_users(where: {_and: {business_id: {_eq: $business_id}, user_type: {_neq: "candidate"}}}) {
                user_type
            }
        }
    }
`

const INVITE_USER = gql`
    mutation ($business_id: Int!, $email: String!, $redirect_link: String!){
        inviteUser(business_id: $business_id, email: $email, redirect_link: $redirect_link) {
            id
            email
        }
    }
`

const GET_USER = gql`
    query GetUser($id: Int!) {
        users_by_pk(id: $id) {
            github_username
            email
            created_at
        }
    }
`
const GET_AUTHED = gql`
    query GetUser($limit: Int = 1) {
        users(limit: $limit) {
            id
            email
            github_username
            auth_id
        }
    }`

export {GET_USERS, INVITE_USER, GET_USER, GET_AUTHED};