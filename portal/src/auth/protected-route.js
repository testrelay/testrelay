import React from "react";
import { Route, Redirect, useLocation } from "react-router-dom";
import { Loading, Sidebar } from "../components/index";
import { useBusiness } from "../components/business/hook";
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

  return (< Route
    component={
      withBusinessRequired(component)
    }
    {...args}
  />)
};

const withBusinessRequired = (Component) => (props) => {
  const location = useLocation();
  const { loading, selected } = useBusiness();
  if (loading) {
    console.log("waiting for master business to load")
    return (
      <Loading />
    )
  }

  console.log("selected in business required", selected)
  if (!selected && !location.pathname.includes("business/create")) {
    return <Redirect push to={{
      pathname: "/business/create",
      state: { referrer: location.pathname, setup: true }
    }} />
  }

  return (
    <Sidebar  {...props} >
      <Component {...props} />
    </Sidebar>
  )
}

export default ProtectedRoute;