import React from "react";
import { Loading } from "../../components";
import { Redirect, useLocation } from 'react-router-dom';
import { useMutation } from "@apollo/client";
import { UPDATE_BUSINESS_GITHUB } from "../components/business/queries";
import { useBusiness } from "../components/business/hook";


const GithubSetup = () => {
  const { loading: businessLoading, selected, setSelected } = useBusiness();
  console.log("selected withint github selected", selected)
  let query = new URLSearchParams(useLocation().search);
  const installationID = query.get("installation_id");
  if (installationID == null) {
    // @todo show error and redirect back to github install
  }

  const [updateBusiness, { loading, data, error }] = useMutation(UPDATE_BUSINESS_GITHUB)

  if (loading || businessLoading) {
    return <Loading />;
  }
  if (error) {
    return <div>{error.message}</div>
  }

  if (data && data.update_businesses_by_pk) {
    setSelected(data.update_businesses_by_pk);

    return <Redirect to="/tests/create" />
  }

  return (
    <div>{updateBusiness({ variables: { id: selected.id, installationId: installationID } })}</div>
  );
}

export default GithubSetup;
