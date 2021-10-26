import React, { useRef, useState } from "react";
import { getFunctions, httpsCallable } from "@firebase/functions";
import { Redirect } from "react-router";
import { Loading } from "../../../components";
import firebase from "../../../auth/firebase";
import { AlertError } from "../../../components/alerts";
import { useBusiness } from "../business/hook";

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
            console.log(email.current.value);
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
            <div className="w-full bg-white p-8 mb-8 shadow-md rounded-xl">
                <div className="flex flex-row space-x-4">
                    <div className="flex-grow">
                        <label className="block uppercase text-gray-700 text-sm font-bold mb-2">
                            User email
                        </label>
                        <input name="name" ref={email} className="input input-bordered w-full text-gray-700" type="email" placeholder="joe@bloggs.com" />
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