import React, {useEffect, useRef, useState} from "react";
import {Redirect} from "react-router";
import {AlertError} from "../../../components/alerts";
import {useBusiness} from "../business/hook";
import {useMutation} from "@apollo/client";
import {INVITE_USER} from "./queries";

const SubmitBtn = (props) => {
    if (props.loading) {
        return (
            <button className="bg-gray-600 text-white text-sm rounded px-4 py-2 w-auto">
                Loading
            </button>
        )
    }

    return (
        <button className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto"
                onClick={props.submit}>
            Invite
        </button>
    )
}

const CreateUser = () => {
    const email = useRef("");
    const {selected} = useBusiness();
    const [error, setError] = useState(null)
    const [redirect, setRedirect] = useState(false)
    const [inviteUser, {data: muData, loading, error: muError}] = useMutation(INVITE_USER);

    useEffect(() => {
        if (muData) {
            setRedirect(true);
        }
    }, [muData]);

    useEffect(() => {
        if (muError) {
            setError("could not invite user " + email.current.value + " please refresh and try again");
        }
    }, [muError])

    const submit = async () => {
        try {
            await inviteUser({
                variables: {
                    business_id: selected.id,
                    email: email.current.value,
                    redirect_link: process.env.REACT_APP_URL + "/assignments/assigned"
                }
            });
        } catch (e) {
        }
    }

    if (redirect) {
        return (
            <Redirect push to="/users"/>
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
                        <input name="name" ref={email} className="input h-9 input-bordered w-full text-gray-700"
                               type="email" placeholder="joe@bloggs.com"/>
                    </div>
                    <div className="flex items-end">
                        <SubmitBtn submit={submit} loading={loading}/>
                    </div>
                </div>
            </div>
            {error &&
            <AlertError message={error}/>
            }
        </div>
    )
}

export default CreateUser;