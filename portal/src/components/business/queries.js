import { gql } from '@apollo/client';

const GET_BUSINESS = gql`
    query GetBusiness {
        businesses {
           id
           name
           github_installation_id
           setup
           creator_id
        }
    }
`

const GET_MASTER_BUSINESS = gql`
query GetMasterBusiness($id: Int!) {
    businesses(limit: 1, where: {creator_id: {_eq: $id}}) {
    id
    name
    setup
    github_installation_id
    creator_id
    }
}
  
`

const UPDATE_BUSINESS_GITHUB = gql`
    mutation UpdateBusiness($id: Int!, $installationId: String!) {
        update_businesses_by_pk(pk_columns: {id: $id}, _set: {github_installation_id: $installationId}) {
            id
            name
            setup
            github_installation_id
            creator_id
        }
    }  
`

const UPDATE_BUSINESS_NAME = gql`
    mutation UpdateBusiness($id: Int!, $name: String!) {
        update_businesses_by_pk(pk_columns: {id: $id}, _set: {name: $name, setup: true}) {
            id
        }
    }  
`

const INSERT_BUSINESS = gql`
    mutation InsertBusiness($name: String!, $user_id: Int!, $user_type: String!) {
        insert_businesses_one(object: {name: $name, business_users: {data: {user_id: $user_id, user_type: $user_type}}}) {
            id
            name
            setup
            github_installation_id
        }
    }
`

export { GET_BUSINESS, UPDATE_BUSINESS_GITHUB, UPDATE_BUSINESS_NAME, INSERT_BUSINESS, GET_MASTER_BUSINESS };