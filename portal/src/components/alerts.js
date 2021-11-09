import React from "react";


const AlertError = (props) => {
    return (
        <div className="alert alert-error">
            <div className="flex-1 item">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="w-5 h-5 mx-2 stroke-current">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
                </svg>
                <label className="text-sm">{props.message}</label>
            </div>
        </div>
    )
}

export { AlertError };