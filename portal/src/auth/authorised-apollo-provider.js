import { ApolloClient, ApolloProvider, createHttpLink, InMemoryCache } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import React from 'react';
import { useFirebaseAuth } from './firebase-hooks';


const AuthorizedApolloProvider = ({ children, }) => {
  const { token } = useFirebaseAuth();
  const httpLink = createHttpLink({
    uri: 'https://api.testrelay.io/v1/graphql'
  })

  const authLink = setContext(async () => {
    return {
      headers: {
        Authorization: `Bearer ${token}`
      }
    };
  });



  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
  });

  return (
    <ApolloProvider client={apolloClient}>
      {children}
    </ApolloProvider>
  );
};

export default AuthorizedApolloProvider;