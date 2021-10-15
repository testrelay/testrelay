import { useMutation } from "@apollo/client";
import React, { useState } from "react";
import { Redirect, useLocation } from "react-router-dom";
import { Loading } from "../../components";
import { AlertError } from "../../components/alerts";
import { useBusiness } from "../../components/business/hook";
import { UPDATE_BUSINESS_NAME } from "../../components/business/queries";

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
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
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
                <div className="w-full bg-white p-8 mb-8 shadow-md rounded-xl">
                    <label className="block uppercase text-gray-700 text-sm font-bold mb-2">
                        Company Name
                    </label>
                    <p className="mb-2">This name will be displayed in emails and correspondance when scheduling assignments with candidates.</p>
                    <input name="name" value={name} onChange={(e) => { setName(e.target.value) }} className="input input-bordered w-full text-gray-700" type="text" placeholder="e.g. BE candidate Test" />
                </div>
                <div className="w-full bg-white p-8 mb-8 shadow-md rounded-xl">
                    <SubmitBtn loading={loading} submit={updateBiz} />
                </div>
            </div>
        </div>
    )
}

export default Settings;