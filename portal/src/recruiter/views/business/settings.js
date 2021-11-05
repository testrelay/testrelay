import { useMutation } from "@apollo/client";
import React, { useState } from "react";
import { Redirect, useLocation } from "react-router-dom";
import { Loading } from "../../../components";
import { AlertError } from "../../../components/alerts";
import { useBusiness } from "../../components/business/hook";
import { UPDATE_BUSINESS_NAME } from "../../components/business/queries";

const SubmitBtn = (props) => {
    if (props.loading) {
        return (
            <button className="bg-gray-500 text-white text-sm rounded px-4 py-2 w-auto">
                Loading
            </button>
        )
    }

    return (
        <button className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto" onClick={props.submit}>
            Save
        </button>
    )
}

const SetupInfo = (props) => {
    if (props.setup) {
        return (
            <div className="alert alert-info mb-6">
                <div className="flex-1">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="w-6 h-6 mx-2 stroke-current">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    <label>Please complete your TestRelay account setup by setting up your company name.</label>
                </div>
            </div>
        )
    }

    return (null);
}

const Settings = () => {
    const location = useLocation();
    location.state = location.state || {};
    const { selected } = useBusiness()
    const [name, setName] = useState(selected.name);
    const [setup, setSetup] = useState(location.state.setup ? location.state.setup : false)
    const [updateName, { loading, error, data }] = useMutation(UPDATE_BUSINESS_NAME);

    const updateBiz = () => {
        updateName({
            variables: { id: selected.id, name }
        }).catch(e => { })
    }

    if (data) {
        if (location.state.referrer) {
            const r = (<Redirect to={location.state.referrer} />)
            delete location.state.referrer;
            delete location.state.setup;

            return r;
        }

        delete location.state.referrer;
        delete location.state.setup;
        setSetup(false);
    }
    return (
        <div>
            <div className="py-4 border-b-4 mb-6">
                <h2 className="text-xl font-bold">Account Settings</h2>
            </div>
            <SetupInfo setup={setup} />
            {error &&
                <AlertError message="could not update company, please try again" />
            }
            <div className="pb-8">
                <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                    <label className="block text-gray-700 font-bold mb-2">
                        Company Name
                    </label>
                    <p className="mb-2 text-sm">This name will be displayed in emails and correspondance when scheduling assignments with candidates.</p>
                    <input name="name" value={name} onChange={(e) => { setName(e.target.value) }} className="input input-bordered w-full text-gray-700" type="text" placeholder="e.g. BE candidate Test" />
                </div>
                <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                    <SubmitBtn loading={loading} submit={updateBiz} />
                </div>
            </div>
        </div>
    )
}

export default Settings;