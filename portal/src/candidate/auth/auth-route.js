import React from "react";
import { Route, Redirect } from "react-router"
import { useLocation } from "react-router-dom";
import { useFirebaseAuth } from "./firebase-hooks";


const AuthedRoute = ({ component, ...args }) => {
    const location = useLocation();
    const { user } = useFirebaseAuth();

    location.state = location.state || {};

    if (user) {
        const path = location.state.referrer || "/assignments";
        delete location.state.referrer;

        return <Redirect to={path} />
    }

    return (<Route component={component} {...args} />)
};

export default AuthedRoute;