import React, { useEffect, useState } from 'react';
import { NetworkStatus, useQuery } from '@apollo/client';
import { Link, useLocation } from 'react-router-dom';
import { useFirebaseAuth } from '../../auth/firebase-hooks';
import { Loading } from '../../../components';
import { GET_ASSIGNED } from '../../components/assignments/queries';
import { getFunctions, httpsCallable } from '@firebase/functions';
import firebase from '../../../auth/firebase';
import { GET_AUTHED } from '../../components/users/queries';
import { AlertError } from '../../../components/alerts';


const Status = (props) => {
    let colour = "orange";

    switch (props.status) {
        case "unsent":
            colour = "gray";
            break;
        case "sending":
            colour = "gray";
            break;
        case "sent":
            colour = "blue";
            break;
        case "scheduled":
            colour = "pink"
            break;
        case "in progress":
            colour = "yellow";
            break;
        case "submitted":
            colour = "green";
            break;
        case "missed":
            colour = "red";
            break;
        default:
            colour = "orange";
    }

    return (
        <span className={"text-" + colour + "-400"}>{props.status}</span>
    )
}

const Grid = () => {
    const location = useLocation();
    const params = new URLSearchParams(location.search);
    const code = params.get('code');

    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const { loading: userLoading, data: userData, refetch, networkStatus } = useQuery(GET_AUTHED);
    useEffect(() => {
        console.log(userLoading, userData, networkStatus)
        const linkGithub = async (data) => {
            if (code != null && data.github_username == null && networkStatus !== NetworkStatus.refetch) {
                setLoading(true);
                const functions = getFunctions(firebase, "europe-west2");
                const authenticateGithub = httpsCallable(functions, "authenticateGithub");

                try {
                    const response = await authenticateGithub({ code });
                    if (response.data.status === "ok") {
                        console.log("refetching")
                        refetch();
                        return;
                    }

                    setError("could not link github to account please try again");
                } catch (e) {
                    setError("could not link github to account please try again");
                }
            }

            setLoading(false);
        }

        if (userLoading || networkStatus === NetworkStatus.refetch) {
            setLoading(true);
        }

        if (userData && !userData.users) {
            setError("could not fetch user information, please refresh the browser")
            setLoading(false);
        }

        if (userData && userData.users.length > 0) {
            setUser(userData.users[0]);
            linkGithub(userData.users[0]);
        }

    }, [code, userLoading, userData, networkStatus, refetch])

    if (loading) {
        return <Loading />
    }

    if (user.github_username == null) {
        return (
            <div className="border-4 border-gray-200 border-dashed p-8 rounded flex justify-center items-center">
                <div className="max-w-md text-center justify-center">
                    {error && <div className="mb-4"><AlertError message={error} /></div>}
                    <p className="text-md mb-4">You'll need to link your account to github before seeing your assigned reviews. This is so we can give you collaborator access to generated repositories.</p>
                    <a href={`https://github.com/login/oauth/authorize?scope=user&client_id=${process.env.REACT_APP_GITHUB_CLIENT_ID}&redirect_uri=${window.location.protocol + '//' + window.location.host + window.location.pathname}`} className="mb-2 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 max-w-xs mx-auto">
                        <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                            <svg className="h-5 w-5 text-gray-100 group-hover:text-gray-200" width="20px" height="20px" viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
                                <g>
                                    <path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" />
                                </g>
                            </svg>
                        </span>
                        Authenticate
                    </a>
                </div>
            </div>
        )
    }

    return (<Assignments />)
}

const Assignments = () => {
    const { user } = useFirebaseAuth();
    const { data, loading } = useQuery(GET_ASSIGNED, { fetchPolicy: "network-only", variables: { uid: user.uid } })

    if (loading) {
        return (<Loading />)
    }

    if (data.assignment_users.length === 0) {
        return (
            <div className="border-4 border-gray-200 border-dashed p-8 rounded flex justify-center items-center">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                </svg>
                <p className="text-center font-medium cursor-pointer text-gray-800">No reviews assigned yet</p>
            </div>
        )
    }

    const assignments = data.assignment_users.map((e, i) => {
        return (
            <div key={i} className="bg-white relative py-4 px-8 text-center md:text-left">
                <div className="grid sm:grid-cols-2 gap-4">
                    <div>
                        <div className="text-md md:text-sm font-medium text-indigo-500 mb-2 capitalize">
                            {e.assignment.test.name}
                        </div>
                        <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            <span className="ml-2">{e.assignment.candidate_name}</span>
                        </div>
                    </div>
                    <div>
                        <div className="mb-1">
                            <span className="text-sm mr-2 font-medium text-gray-500">Assignment <Status status={e.assignment.status} /></span>
                        </div>
                        <div>
                            {e.assignment.github_repo_url
                                ? <a className="text-sm text-indigo-500" href={e.assignment.github_repo_url} rel="noreferrer" target="_blank">{e.assignment.github_repo_url}</a>
                                : <span className="text-sm text-indigo-300">repo yet to be generated</span>
                            }
                        </div>
                    </div>
                </div>
            </div>
        )

    })

    return (<div className="bg-white overflow-hidden rounded-lg shadow-md divide-y-2">
        {assignments}
    </div>)
}

const Assigned = () => {
    return (
        <div>
            <div className="py-2 border-b-4 mb-6 flex">
                <div className="flex-1 flex flex-col justify-center">
                    <h2 className="text-xl font-bold">Assigned Reviews</h2>
                </div>
                <div>
                    <Link to="/assignments" className="flex items-center text-indigo-500 hover:text-gray-500">
                        <span>see all</span>
                        <svg xmlns="http://www.w3.org/2000/svg" className="ml-1 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </Link>
                </div>
            </div>
            <Grid />
        </div>
    )
}

export default Assigned;