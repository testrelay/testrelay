import { useMutation, useQuery } from "@apollo/client";
import React, { useEffect, useState } from "react";
import { Redirect, useLocation } from "react-router-dom";
import { Loading } from "../../../components";
import { AlertError } from "../../../components/alerts";
import { GET_BUSINESS, INSERT_BUSINESS } from "../../components/business/queries";
import { getFunctions, httpsCallable } from "firebase/functions";
import firebase from "../../../auth/firebase";
import { useFirebaseAuth } from "../../auth/firebase-hooks";
import { useBusiness } from "../../components/business/hook";

const SubmitBtn = (props) => {
    if (props.loading) {
        return (
            <button className="btn btn-disabled">
                <Loading />
            </button>
        )
    }

    return (
        <button className="btn btn-primary bg-indigo-600" onClick={props.submit}>
            Save Organisation
        </button>
    )
}

const SetupInfo = (props) => {
    if (props.setup) {
        return (
            <div className="alert alert-info mb-6">
                <div className="flex-1">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="w-6 h-6 mx-2 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    <label>Please complete your TestRelay account setup by setting up your organisation.</label>
                </div>
            </div>
        )
    }

    return (null);
}

const Create = (props) => {
    const location = useLocation();
    location.state = location.state || {};
    const [name, setName] = useState("")
    const [r, setRedirect] = useState(false);
    const { user, claims } = useFirebaseAuth();
    const { setSelected } = useBusiness();

    const [pageLoading, setPageLoading] = useState(true);
    const [loading, setLoading] = useState(false);

    const [error, setError] = useState(null);

    const { data } = useQuery(GET_BUSINESS, { fetchPolicy: "network-only" })

    const [insertBusiness, { error: queryError }] = useMutation(INSERT_BUSINESS, {
        onCompleted: async (data) => {
            console.log("update meta with business id", data.insert_businesses_one.id);
            const functions = getFunctions(firebase, "europe-west2");
            const changeMeta = httpsCallable(functions, "changeMeta");
            await changeMeta({ user_type: "recruiter", business_id: data.insert_businesses_one.id });
            console.log("refreshing the id token after appending business");
            await user.getIdToken(true);
            setSelected(data.insert_businesses_one);
            setLoading(false);
            setRedirect(location.state.referrer || "/tests");
        }
    });

    useEffect(() => {
        if (queryError) {
            setError(queryError);
            setLoading(false);
        }
    }, [queryError])

    useEffect(() => {
        if (pageLoading) {
            setPageLoading(true);
        }
    }, [pageLoading]);

    useEffect(() => {
        if (data) {
            if (data.businesses.length > 0) {
                const hasBusiness = data.businesses.find((e) => {
                    return e.creator_id === parseInt(claims["x-hasura-user-pk"]);
                })

                if (hasBusiness) {
                    return <Redirect to="/tests" />
                }
            }

            setPageLoading(false)
        }
    }, [data, claims])


    const insert = async () => {
        setLoading(true);
        console.log("inserting with user id " + claims["x-hasura-user-pk"]);
        insertBusiness({
            variables: { name, user_id: parseInt(claims["x-hasura-user-pk"]), user_type: "recruiter" },
        }).catch(e => {
            setError(e)
        });
    }


    if (r) {
        console.log("redirecting from business create to: ", r);
        delete location.state.referrer;
        delete location.state.setup;
        return (<Redirect to={r} />)
    }

    if (error) {
        console.error(error);
    }

    if (pageLoading) {
        return <Loading />
    }

    return (
        <div>
            <div className="py-4 border-b-4 mb-6">
                <h2 className="text-xl font-bold">Create your organisation</h2>
            </div>
            <SetupInfo setup={location.state.setup} />
            <div className="pb-8">
                <div className="w-full bg-white p-8 mb-8 shadow-md rounded-xl">
                    <label className="block uppercase text-gray-700 text-sm font-bold mb-2">
                        Organisation Name
                    </label>
                    <p className="mb-2">This name will be displayed in emails and correspondance when scheduling assignments with candidates.</p>
                    <input name="name" value={name} onChange={(e) => { setName(e.target.value) }} className="input input-bordered w-full text-gray-700" type="text" placeholder="e.g. BE candidate Test" />
                </div>
                <div className="w-full bg-white p-8 mb-8 shadow-md rounded-xl">
                    <SubmitBtn loading={loading} submit={insert} />
                </div>
                {error &&
                    <AlertError message="could not create organisation, please try again" />
                }
            </div>
        </div>
    )
}

export default Create;