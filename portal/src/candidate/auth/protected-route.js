import React from "react";
import {Redirect, Route, useLocation} from "react-router-dom";
import {Loading} from "../components";
import {useFirebaseAuth} from "../../auth/firebase-hooks";

const ProtectedRoute = ({component, ...args}) => {
    const location = useLocation()
    const {user, loading} = useFirebaseAuth();

    if (loading) {
        return (<Loading/>)
    }

    if (!user) {
        return (<Redirect push to={{
            pathname: "/login",
            state: {referrer: location.pathname}
        }}/>)
    }

    return (< Route component={component} {...args} />)
}

export default ProtectedRoute;
