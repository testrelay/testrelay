import React from "react";
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";

import {  GithubSetup } from "./views";
import { AssignmentView, CandidateListView } from "./views/assignments";
import { TestListView, CreateTest } from "./views/tests";
import ProtectedRoute from "./auth/protected-route";

import CandidateCreate from "./views/assignments/create";
import TestView from "./views/tests/view";

import "../app.css";
import Settings from "./views/business/settings";
import Login from "./views/login";
import Register from "./views/register";
import AuthedRoute from "./auth/auth-route";
import UserList from "./views/users/list";
import AuthorizedApolloProvider from "./auth/authorised-apollo-provider";
import { FirebaseAuthProvider } from "./auth/firebase-hooks";
import Create from "./views/business/create";
import { BusinessProvider } from "./components/business/hook";
import UserCreate from "./views/users/view";
import View from "./views/users/view";
import Assigned from "./views/assignments/assigned";

const App = () => {

  return (
    <div id="app" className="h-full">
      <BrowserRouter>
        <FirebaseAuthProvider>
          <AuthorizedApolloProvider >
            <BusinessProvider>
              <Switch>
                <Route exact path="/" >
                  <Redirect to="/login" />
                </Route>


                <AuthedRoute path="/login" exact component={Login} />
                <AuthedRoute path="/register" exact component={Register} />

                <ProtectedRoute path="/github-setup" exact component={GithubSetup} />

                <ProtectedRoute path="/assignments" exact component={CandidateListView} />
                <ProtectedRoute path="/assignments/create" exact component={CandidateCreate} />
                <ProtectedRoute path="/assignments/assigned" exact component={Assigned} />
                <ProtectedRoute path="/assignments/:id/view" exact component={AssignmentView} />

                <ProtectedRoute path="/tests" exact component={TestListView} />
                <ProtectedRoute path="/tests/create" exact component={CreateTest} />
                <ProtectedRoute path="/tests/:id/view" exact component={TestView} />

                <ProtectedRoute path="/settings" exact component={Settings} />
                <ProtectedRoute path="/business/create" exact component={Create} />

                <ProtectedRoute path="/users" exact component={UserList} />
                <ProtectedRoute path="/users/create" exact component={UserCreate} />
                <ProtectedRoute path="/users/:id/view" exact component={View} />
              </Switch>
            </BusinessProvider>
          </AuthorizedApolloProvider>
        </FirebaseAuthProvider>
      </BrowserRouter>
    </div>
  );
};

export default App;
