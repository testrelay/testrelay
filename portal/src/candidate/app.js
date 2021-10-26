import React from "react";
import { BrowserRouter, Switch } from "react-router-dom";
import ProtectedRoute from "./auth/protected-route";

import AssignmentView from "./views/assignments/view";
import List from "./views/assignments/list";

import "../app.css";
import Login from "./views/login";
import AuthedRoute from "./auth/auth-route";
import AuthorizedApolloProvider from "./auth/authorised-apollo-provider";
import { FirebaseAuthProvider } from "./auth/firebase-hooks";

const App = () => {
  return (
    <div id="app" className="h-full">
      <BrowserRouter>
        <FirebaseAuthProvider>
          <AuthorizedApolloProvider>
            <Switch>
              <AuthedRoute path="/login" exact component={Login} />

              <ProtectedRoute path="/assignments/:id/view" exact component={AssignmentView} />
              <ProtectedRoute path="/assignments" exact component={List} />
            </Switch>
          </AuthorizedApolloProvider>
        </FirebaseAuthProvider>
      </BrowserRouter>
    </div>
  );
};

export default App;
