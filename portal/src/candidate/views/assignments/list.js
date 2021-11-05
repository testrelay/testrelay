import React, {useState} from "react";

import {useMutation, useQuery} from "@apollo/client";
import {GET_ASSIGNMENTS, UPDATE_ASSIGNMENT_CANCELED} from "../../components/assignments/queries";
import {Loading} from "../../components";
import {ErrorAlert} from "../../components/alert";
import {assignmentLimit} from "../../../components/time";
import {Link} from "react-router-dom";
import {useFirebaseAuth} from "../../auth/firebase-hooks";

const Buttons = (props) => {
    const [cancel, {loading, error}] = useMutation(UPDATE_ASSIGNMENT_CANCELED, {variables: {id: props.id}});
    const [e, setError] = useState(null);

    if (props.test_day_chosen == null) {
        return (
            <div className="flex flex-wrap items-start space-x-2 flex-row justify-center md:justify-end">
                <Link className="btn btn-primary" to={"/assignments/" + props.id + "/view"}>schedule</Link>
            </div>
        )
    }

    if (loading) {
        return (
            <div className="flex flex-wrap items-start space-x-2 flex-row justify-center md:justify-end">
                <button className="btn btn-disabled">reschedule</button>
                <button className="btn btn-disabled"><Loading/></button>
            </div>
        )
    }


    const click = async () => {
        try {
            await cancel();
        } catch (error) {
            setError(error)
        }
    }

    const isError = () => {
        return e || error;
    }
    return (
        <div className="flex flex-wrap items-start space-x-2 flex-row justify-center md:justify-end">
            <Link className="btn btn-primary h-10 min-h-0 shadow-lg"
                  to={"/assignments/" + props.id + "/view"}>reschedule</Link>
            <button className="btn btn-warning h-10 min-h-0 shadow-lg" onClick={click}>cancel</button>
            {isError() &&
            <div className="mt-2"><ErrorAlert message="could not cancel assignment, please try again"/></div>}
        </div>
    )
}

const Schedule = (props) => {
    const day = new Date(props.test_day_chosen);
    const time = props.test_time_chosen;
    const choose_until = new Date(props.choose_until);
    const monthNames = ["January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ];

    if (day == null) {
        const text = "Unscheduled, choose until " + monthNames[choose_until.getMonth()] + " " + choose_until.getDate() + ", " + choose_until.getFullYear();

        return (
            <div className="text-sm font-medium text-error">
                {text}
            </div>
        )
    }

    const text = "Test scheduled for " + monthNames[day.getMonth()] + " " + day.getDate() + ", " + day.getFullYear();

    return (
        <div>
            <div className="text-sm font-medium text-gray-900">
                {text}
            </div>
            <div className="text-sm text-gray-500">
                at {time}
            </div>
        </div>
    )
}

const Assignments = (props) => {
    return props.assignments.map((e, i) => {
        if (e.status === "cancelled") {
            return (<div key={i} className="relative bg-white shadow-md p-8 rounded-lg text-center md:text-left">
                    <div className="bg-gray-800 opacity-5 absolute w-full h-full top-0 left-0 right-0 z-10"/>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div>
                            <div className="text-md md:text-sm font-medium text-primary mb-2">
                                {e.test.business.name}
                            </div>
                            <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 flex-shrink-0" fill="none"
                                     viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                </svg>
                                <span className="ml-1">{assignmentLimit(e.time_limit)}</span>
                            </div>
                        </div>
                        <div>
                            <p className="text-warning">Cancelled</p>
                        </div>
                        <div/>
                    </div>
                </div>
            )
        }
        return (
            <div key={i} className="relative shadow-md bg-white p-8 rounded-lg text-center md:text-left">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                        <div className="text-md md:text-sm font-medium text-primary mb-2">
                            {e.test.business.name}
                        </div>
                        <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 flex-shrink-0" fill="none"
                                 viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                            </svg>
                            <span className="ml-1">{assignmentLimit(e.time_limit)}</span>
                        </div>
                    </div>
                    <div>
                        <Schedule {...e} />
                    </div>
                    <div>
                        <Buttons test_day_chosen={e.test_time_chosen} id={e.id}/>
                    </div>
                </div>
            </div>
        )
    })
}

const List = () => {
    const {loading: claimLoading, user} = useFirebaseAuth();
    const {loading, error, data} = useQuery(GET_ASSIGNMENTS, {
            fetchPolicy: "network-only"
        }
    );
    console.log(user);

    if (loading || claimLoading) {
        return (
            <div className="container mx-auto px-4 max-w-2xl">
                <div className="mt-14">
                    <Loading/>
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="container mx-auto px-4 max-w-2xl">
                <div className="mt-14">
                    <ErrorAlert message="could not display profile, please refresh browser and try again"/>
                </div>
            </div>
        )
    }

    if (data) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
                <div className="max-w-4xl w-full">
                    <div className="mb-4 text-center">
                        <p className="text-xl mb-2">Hey {user.displayName.split(" ")[0]}</p>
                        <p className="text-md text-gray-500">Your assignments:</p>
                    </div>
                    <Assignments assignments={data.assignments}/>
                </div>
            </div>
        )
    }

    return (
        <Loading/>
    )
}

export default List;