import React, { useRef, useState } from "react";
import { getFunctions, httpsCallable } from "@firebase/functions";
import { Redirect } from "react-router";
import firebase from "../../../auth/firebase";
import { AlertError } from "../../../components/alerts";
import { useBusiness } from "../business/hook";

const SubmitBtn = (props) => {
    if (props.loading) {
        return (
            <button className="bg-gray-600 text-white text-sm rounded px-4 py-2 w-auto">
                Loading
            </button>
        )
    }

    return (
        <button className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto" onClick={props.submit}>
            Invite
        </button>
    )
}

const CreateUser = () => {
    const email = useRef("");
    const { selected } = useBusiness();
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)
    const [redirect, setRedirect] = useState(false)

    const submit = async () => {
        setLoading(true);

        const functions = getFunctions(firebase, "europe-west2");
        const invite = httpsCallable(functions, "inviteUser");

        try {
            await invite({ email: email.current.value, business_name: selected.name, business_id: selected.id });
            setRedirect(true);
            setLoading(false);
        } catch (error) {
            setError("could not invite user " + email.current.value + " please refresh and try again");
            setLoading(false);
        }
    }

    if (redirect) {
        return (
            <Redirect push to="/users" />
        )
    }

    return (
        <div className="pb-8">
            <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                <div className="flex flex-row space-x-4">
                    <div className="flex-grow">
                        <label className="block text-gray-700 text-sm font-bold mb-2">
                            User email
                        </label>
                        <input name="name" ref={email} className="input h-9 input-bordered w-full text-gray-700" type="email" placeholder="joe@bloggs.com" />
                    </div>
                    <div className="flex items-end">
                        <SubmitBtn submit={submit} loading={loading} />
                    </div>
                </div>
            </div>
            {error &&
                <AlertError message={error} />
            }
        </div>
    )
}

export default CreateUser;