import { useQuery } from "@apollo/client";
import { getAuth, signOut } from "firebase/auth";
import React from "react";
import { Link, useLocation } from "react-router-dom";
import firebase from "../auth/firebase";
import { useFirebaseAuth } from "../recruiter/auth/firebase-hooks";
import { useBusiness } from "../recruiter/components/business/hook";
import { GET_BUSINESS } from "../recruiter/components/business/queries";
import Loading from "./loading";


const BusinessSelect = (props) => {
    const { loading, data } = useQuery(GET_BUSINESS);
    const { claims } = useFirebaseAuth();

    if (loading) {
        return (<Loading />)
    }

    let hasOwnOrg = false;

    if (data.businesses.length === 1) {
        if (data.businesses[0].creator_id === parseInt(claims["x-hasura-user-pk"])) {
            hasOwnOrg = true;
        }

        const name = (
            <div className="flex justify-center items-center p-4">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <p className="text-white ml-2">{data.businesses[0].name}</p>
            </div>
        )

        if (hasOwnOrg) {
            return name
        }

        return (
            <div className="flex justify-center space-x-1 items-center">
                {name}
                <Link to="/business/create" className="btn btn-primary bg-indigo-600 h-10 min-h-0">Create Org</Link>
            </div>
        )
    }

    if (data.businesses.length === 0) {
        return (
            <div className="flex justify-center space-x-1 items-center p-2">
                <Link to="/business/create" className="btn btn-primary bg-indigo-600 h-10 min-h-0">Create Org</Link>
            </div>
        )

    }

    const businesses = data.businesses.map((e) => {
        if (e.creator_id === parseInt(claims["x-hasura-user-pk"])) {
            hasOwnOrg = true;
        }

        return (<option key={e.id} value={e.id}>{e.name}</option>)
    })

    const change = (ev) => {
        const selected = data.businesses.find(e => e.id === ev.target.value);
        props.setSelected(selected);
    }

    const select = (
        <div className="flex justify-center items-center p-2">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <select onChange={change} className="select select-bordered ml-2 min-h-0 h-10" value={props.selected.id} defaultValue={props.selected.id}>
                {businesses}
            </select>
        </div>
    )

    if (!hasOwnOrg) {
        return (
            <div className="flex justify-center space-x-1 items-center">
                {select}
                <Link to="/business/create" className="btn btn-primary bg-indigo-600 h-10 min-h-0">Create Org</Link>
            </div>
        )
    }


    return select
}
const Sidebar = (props) => {
    const location = useLocation();
    const { selected, setSelected } = useBusiness();
    const pieces = location.pathname.split("/");


    let path = "assignments";
    if (pieces.length > 1) {
        path = pieces[1];
    }

    const isSelected = (link) => {
        if (link === path) {
            return "text-gray-800";
        }

        return "text-indigo-500";
    }

    const assignments = () => {
        if ("assignments" === path) {
            return "px-8 py-2 block";
        }

        return "hidden";
    }

    const revoke = async () => {
        const auth = getAuth(firebase);


        setSelected(null);
        await signOut(auth)
    }

    return (

        <div className="bg-gray-100 drawer drawer-mobile h-full">
            <input id="my-drawer-2" type="checkbox" className="drawer-toggle" />
            <div className="flex flex-col drawer-content">
                <label htmlFor="my-drawer-2" className="mb-4 btn btn-primary bg-indigo-600 drawer-button lg:hidden">open menu</label>
                <div className="bg-gray-800 p-2 shadow-md clear-both">
                    <div className="float-right">
                        <BusinessSelect selected={selected} setSelected={setSelected} />
                    </div>
                </div>
                <div className="container mx-auto py-10 h-64 md:w-4/5 w-11/12 px-6">
                    {props.children}
                </div>
            </div>
            <div className="drawer-side gap-1 shadow-xl bg-white">
                <label htmlFor="my-drawer-2" className="drawer-overlay"></label>
                <ul className="menu py-4 overflow-y-auto w-80 text-base-content">
                    <li className="mb-6 px-4">
                        <div className="flex pb-6 pl-4 items-center border-b-4">
                            <div className="flex items-center bg-gray-800 shadow w-14 h-14 rounded-full">
                                <h1 className="mx-auto font-semibold text-lg text-white">
                                    <svg height="30pt" viewBox="0 -48 480 480" width="30pt" fill="#fff" xmlns="http://www.w3.org/2000/svg"><path d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0" /><path d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0" /></svg>
                                </h1>
                            </div>
                            <h1 className="ml-2 text-xl font-extrabold">TestRelay</h1>
                        </div>
                    </li>
                    <li className={isSelected("tests") + " hover:text-gray-500 cursor-pointer"}>
                        <Link to="/tests">
                            <div className="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                </svg>
                                <span className="text-sm ml-2">Tests</span>
                            </div>
                        </Link>
                    </li>
                    <li className={isSelected("assignments") + " cursor-pointer"}>
                        <Link className="hover:text-gray-500" to="/assignments">
                            <div className="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path d="M12 14l9-5-9-5-9 5 9 5z" />
                                    <path d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222" />
                                </svg>
                                <span className="text-sm  ml-2">Assignments</span>
                            </div>
                        </Link>
                        <div className={assignments()}>
                            <Link className="flex items-center text-gray-800 hover:text-indigo-500" to="/assignments/assigned">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                <span className="text-xs">To review</span>
                            </Link>
                        </div>
                    </li>
                    <li className={isSelected("users") + " hover:text-gray-500 cursor-pointer"}>
                        <Link to="/users">
                            <div className="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                                </svg>
                                <span className="text-sm  ml-2">Users</span>
                            </div>
                        </Link>
                    </li>
                    <li className={isSelected("settings") + " hover:text-gray-500 cursor-pointer"}>
                        <Link to="/settings">
                            <div className="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                </svg>
                                <span className="text-sm  ml-2">Settings</span>
                            </div>
                        </Link>
                    </li>
                </ul>
                <div className="flex px-4 mb-4 items-end">
                    <button className="btn btn-primary bg-indigo-600 btn-sm  rounded-md" onClick={revoke}>Logout</button>
                </div>
            </div>
        </div>
    );
}

export default Sidebar;
