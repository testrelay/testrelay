import {ApolloClient, ApolloProvider, from, fromPromise, HttpLink, InMemoryCache} from '@apollo/client';
import {setContext} from '@apollo/client/link/context';
import React from 'react';
import {useFirebaseAuth} from './firebase-hooks';
import {onError} from "@apollo/client/link/error";


const AuthorizedApolloProvider = ({role, children}) => {
    const {user, token} = useFirebaseAuth();
    const httpLink = new HttpLink({
        uri: process.env.REACT_APP_GRAPHQL_URL,

    })

    const authLink = setContext(() => {
        const ctx = {
            headers: {
                Authorization: `Bearer ${token}`
            }
        };

        if (role != null) {
            ctx.headers['X-Hasura-Role'] = role;
        }

        return ctx;
    });

    const errorLink = onError(({graphQLErrors, forward, operation}) => {
            if (graphQLErrors) {
                for (const error of graphQLErrors) {
                    if (error.message.includes("JWTExpired")) {
                        return fromPromise(user.getIdToken(true)).filter((value) => Boolean(value)).flatMap(newToken => {
                            const oldHeaders = operation.getContext().headers;
                            operation.setContext({
                                headers: {
                                    ...oldHeaders,
                                    Authorization: `Bearer ${newToken}`
                                }
                            })

                            return forward(operation);
                        })
                    }

                }
            }
        }
    );

    const apolloClient = new ApolloClient({
        link: from([
            authLink,
            errorLink,
            httpLink,
        ]),
        cache: new InMemoryCache(),
    });

    return (
        <ApolloProvider client={apolloClient}>
            {children}
        </ApolloProvider>
    );
};

export default AuthorizedApolloProvider;