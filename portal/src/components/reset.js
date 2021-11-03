import React, {useRef, useState} from "react";
import {getAuth, sendPasswordResetEmail} from "firebase/auth";
import {Link} from "react-router-dom";

const Reset = (props) => {
    const [message, setMessage] = useState()
    const email = useRef();

    const reset = () => {
        if (email.current.value === "") {
            setMessage({type: "error", msg: "Please provide a email"})
            return
        }
        sendPasswordResetEmail(getAuth(), email.current.value)
            .then(() => {
                setMessage({type: "success", msg: "Check your inbox for a password reset email"})
            })
            .catch((error) => {
                console.log(error);
                setMessage({type: "error", msg: "Failed to process password reset request, please try again"})
            });
    }
    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full ">
                <div className="bg-white p-8 shadow-md rounded-lg">
                    <h2 className="mb-4 text-center text-3xl font-bold text-gray-900">
                        Reset your password
                    </h2>
                    <input name="email" type="email" autoComplete="email" required
                           className="mb-4 rounded appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                           placeholder="Email address" ref={email}/>

                    <button onClick={reset}
                            className="mb-2 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                        Reset
                    </button>
                    {message &&
                    <div className={"mt-4 alert alert-" + message.type}>
                        <label className="text-sm">{message.msg}</label>
                    </div>
                    }
                </div>
                <div className="text-right">
                    <Link to="/login" className="flex justify-center items-center p-4">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-indigo-500" fill="none"
                             viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M11 15l-3-3m0 0l3-3m-3 3h8M3 12a9 9 0 1118 0 9 9 0 01-18 0z"/>
                        </svg>
                        <span className="text-indigo-500 ml-1">
                            login
                        </span>
                    </Link>
                </div>
            </div>
        </div>
    )
}

export default Reset;
