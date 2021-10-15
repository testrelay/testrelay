import React from "react";
import { Route, Redirect, useLocation } from "react-router-dom";
import { Loading } from "../components/index";
import { useFirebaseAuth } from "./firebase-hooks";

const ProtectedRoute = ({ component, ...args }) => {
  const location = useLocation()
  const { user, loading } = useFirebaseAuth();

  if (loading) {
    console.log("waiting for auth to load")
    return (<Loading />)
  }

  if (!user) {
    return (<Redirect push to={{
      pathname: "/login",
      state: { referrer: location.pathname }
    }} />)
  }

  return (< Route component={component} {...args} />)
}

export default ProtectedRoute;
