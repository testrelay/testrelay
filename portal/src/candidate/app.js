import React from "react";
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";
import ProtectedRoute from "./auth/protected-route";

import AssignmentView from "./views/assignments/view";
import List from "./views/assignments/list";

import "../app.css";
import Login from "./views/login";
import AuthedRoute from "../auth/auth-route";
import AuthorizedApolloProvider from "../auth/authorised-apollo-provider";
import {CandidateAuthProvider} from "../auth/firebase-hooks";
import ResetView from "./views/reset";

const App = () => {
    return (
        <div id="app" className="h-full">
            <BrowserRouter>
                <CandidateAuthProvider>
                    <AuthorizedApolloProvider role="candidate">
                        <Switch>
                            <Route exact path="/"><Redirect to="/login"/></Route>

                            <AuthedRoute redirect="/assignments" path="/login" exact component={Login}/>
                            <AuthedRoute redirect="/assignments" path="/password-reset" exact component={ResetView}/>

                            <ProtectedRoute path="/assignments/:id/view" exact component={AssignmentView}/>
                            <ProtectedRoute path="/assignments" exact component={List}/>
                        </Switch>
                    </AuthorizedApolloProvider>
                </CandidateAuthProvider>
            </BrowserRouter>
        </div>
    );
};

export default App;
