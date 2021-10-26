import React from "react";
import { Route, Redirect } from "react-router";
import { useLocation } from "react-router-dom";
import { Loading } from "../../components";
import { useFirebaseAuth } from "./firebase-hooks";


const AuthedRoute = ({ component, ...args }) => {
    const location = useLocation();
    const { user, loading } = useFirebaseAuth();

    location.state = location.state || {};

    if (loading) {
        return <Loading />
    }

    if (user) {
        const path = location.referrer || "/tests";
        delete location.state.referrer;

        return <Redirect to={path} />
    }

    return (<Route component={component} {...args} />)
};

export default AuthedRoute;