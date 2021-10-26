import React, { useEffect, useState } from "react";
import DatePicker from 'react-datepicker';
import moment from "moment-timezone";
import { useMutation, useQuery } from "@apollo/client";
import { GET_ASSIGNMENT, UPDATE_ASSIGNMENT_WITH_SCHEDULE } from "../../components/assignments/queries";
import { Redirect, useLocation, useParams } from "react-router-dom";
import { Loading } from "../../components";
import { ErrorAlert } from "../../components/alert";
import { assignmentLimit } from "../../components/assignments/time";
import { GET_USER } from "../../components/queries";
import firebase from "../../../auth/firebase";
import { getFunctions, httpsCallable } from "firebase/functions";


const TimezoneSelect = (props) => {
    const getTimeZoneList = moment.tz.names()
        .map(t => {
            const z = moment.tz(t).format("Z");
            return {
                name: `(GMT${z}) ${t}`,
                value: t,
                time: z
            }
        });

    const sortByZone = (a, b) => {
        let [ahh, amm] = a.time.split(":");
        let [bhh, bmm] = b.time.split(":");
        return (+ahh * 60 + amm) - (+bhh * 60 + bmm)
    };

    const options = getTimeZoneList.sort(sortByZone).map((e, i) => {
        return (<option key={i} value={e.value}>{e.name}</option>)
    })

    return (
        <select name="test_timezone_chosen" className="select select-bordered w-full max-w-xs" onChange={props.change} value={props.value}>
            {options}
        </select>
    )
}

const TimeSelect = (props) => {
    const hours = Array.from({
        length: 48
    }, (_, hour) => moment({
        hour: Math.floor(hour / 2),
        minutes: (hour % 2 === 0 ? 0 : 30)
    }).format('HH:mm')).map((e, i) => {
        return <option key={i} value={e}>{e}</option>
    });

    return (
        <select name="test_time_chosen" className="select select-bordered w-full max-w-xs" onChange={props.change} value={props.value}>
            {hours}
        </select>
    )
}

const SubmitButton = (props) => {
    if (props.loading) {
        return (
            <button className="btn btn-disabled"><Loading /></button>
        )
    }

    return (<button className="btn btn-primary" onClick={props.click}>Submit</button>);
}

const TestBody = ({ code, assignment }) => {
    const defaultTimeZone = moment.tz.guess();
    const [form, setForm] = useState({
        test_day_chosen: new Date(),
        test_time_chosen: "00:00",
        test_timezone_chosen: defaultTimeZone
    });

    const endDate = new Date(assignment.choose_until);

    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [formLoading, setFormLoading] = useState(false);
    const [error, setError] = useState(null);
    const [redirect, setRedirect] = useState(false);
    const { loading: userLoading, data: userData, refetch } = useQuery(GET_USER, { fetchPolicy: 'network-only' });

    const [
        updateAssignment,
        updateAssignmentRes,
    ] = useMutation(UPDATE_ASSIGNMENT_WITH_SCHEDULE);

    useEffect(() => {
        if (assignment) {
            if (assignment.test_day_chosen != null) {
                const split = assignment.test_time_chosen.split(":");

                setForm(f => {
                    return {
                        ...f,
                        test_timezone_chosen: assignment.test_timezone_chosen,
                        test_time_chosen: split[0] + ":" + split[1],
                        test_day_chosen: new Date(assignment.test_day_chosen)
                    }
                })
            }
        }
    }, [assignment]);

    useEffect(() => {
        const linkGithub = async (data) => {
            if (code != null && data.github_username == null) {
                console.log("linking github with data", data)
                setLoading(true);
                const functions = getFunctions(firebase, "europe-west2");
                const authenticateGithub = httpsCallable(functions, "authenticateGithub");

                try {
                    const response = await authenticateGithub({ code });
                    if (response.data.status === "ok") {
                        console.log("refetching")
                        refetch();
                        return;
                    } else {
                        setError("could not link github to account please try again");
                    }
                } catch (e) {
                    setError("could not link github to account please try again");
                }
            }

            setLoading(false);
        }

        if (userLoading) {
            setLoading(true);
        }

        if (userData && !userData.users) {
            setError("could not fetch user information, please refresh the browser")
            setLoading(false);
        }

        if (userData && userData.users.length > 0) {
            console.log("setting user information");
            setUser(userData.users[0]);
            linkGithub(userData.users[0]);
        }

    }, [code, userLoading, userData, refetch])


    useEffect(() => {
        if (updateAssignmentRes.error) {
            setError(updateAssignment.error);
        }

        if (updateAssignment.data) {
            setRedirect(true);
        }

        if (updateAssignment.loading) {
            setFormLoading(updateAssignment.loading);
        }
    }, [updateAssignmentRes, updateAssignment]);

    const inputChange = (e) => {
        setForm(f => { return { ...f, [e.target.name]: e.target.value } });
    }

    const dateChange = (date) => {
        setForm(f => { return { ...f, test_day_chosen: date } });
    }

    const submitForm = () => {
        const meta = Object.assign({}, form, { test_day_chosen: formatDate(form.test_day_chosen) });
        const variables = Object.assign({}, meta, { meta, id: assignment.id })
        updateAssignment({ variables }).catch(e => {
            setError("could not set assignment choices please refresh the browser and try again");
        })

        setRedirect(true);
    }

    if (redirect) {
        return (
            <Redirect push to="/assignments" />
        )
    }

    if (loading) {
        return (<Loading />)
    }

    if (user && user.github_username) {
        return (
            <div>
                <div className="mb-4">
                    <p className="text-xl text-primary mb-4">Schedule a day & time for you to take the test.</p>
                    <label className="label"><span className="label-text">Choose a day to take your test</span></label>
                    <div className="relative">
                        <DatePicker
                            selected={form.test_day_chosen}
                            onChange={(date) => dateChange(date)}
                            maxDate={endDate}
                            minDate={new Date()}
                            nextMonthButtonLabel=">"
                            previousMonthButtonLabel="<"
                        />
                    </div>
                </div>
                <div className="mb-4">
                    <label className="label"><span className="label-text">Choose time you want to start the test</span></label>
                    <TimeSelect change={inputChange} value={form.test_time_chosen} />
                </div>
                <div className="mb-4">
                    <label className="label"><span className="label-text">And what timezone you'll be in when taking the test.</span></label>
                    <TimezoneSelect defaultTimeZone={defaultTimeZone} change={inputChange} value={form.test_timezone_chosen} />
                </div>
                <SubmitButton click={submitForm} loading={formLoading} />
                {error && <ErrorAlert message={error} />}
            </div>
        )
    }

    return (
        <div className="border-t-2 pt-4">
            {error && <div className="mb-4"><ErrorAlert message={error} /></div>}
            <p className="text-md mb-4">Before scheduling your technical test, you'll need to link your account to github. This is needed so TestRelay can invite you to the private github repo where you'll take your test.</p>
            <a href={`https://github.com/login/oauth/authorize?scope=user&client_id=${process.env.REACT_APP_GITHUB_CLIENT_ID}&redirect_uri=${window.location.protocol + '//' + window.location.host + window.location.pathname}`} className="mb-2 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                    <svg className="h-5 w-5 text-gray-100 group-hover:text-gray-200" width="20px" height="20px" viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
                        <g>
                            <path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" ></path>
                        </g>
                    </svg>
                </span>
                Authenticate with github
            </a>
        </div>
    )

}

const AssignmentView = () => {
    const location = useLocation();
    const params = new URLSearchParams(location.search);
    const code = params.get('code');

    const id = useParams().id;

    const { loading, error, data } = useQuery(GET_ASSIGNMENT, {
        variables: { id },
        fetchPolicy: 'network-only',
    });

    if (loading) {
        return (
            <Loading />
        );
    }

    if (error) {
        console.log("error from fetching assignment", error)
        return (
            <div className="container mx-auto px-4 max-w-2xl">
                <div className="mt-14">
                    <ErrorAlert message="failed to fetch assignment, please reload the page" />
                </div>
            </div>
        )
    }

    if (data.assignments_by_pk.time_limit == null) {
        return (
            <div className="container mx-auto px-4 max-w-2xl">
                <div className="mt-14">
                    <ErrorAlert message="this assignment does not exist, check you have the correct permissions to view it" />
                </div>
            </div>
        )
    }


    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-lg w-full p-8 shadow-lg rounded-lg bg-white">
                <div className="mb-8">
                    <h1 className="text-2xl mb-4">Hey Hugo</h1>
                    <p className="text-md"><b>{data.assignments_by_pk.test.business.name}</b> has invited you to take a technical test.<br />
                        You will have <b>{assignmentLimit(data.assignments_by_pk.time_limit)}</b> to take the test and complete it wil one of the following programming languages: <b>{languages(data.assignments_by_pk.test.test_languages)}</b></p>
                </div>
                <TestBody
                    code={code}
                    assignment={data.assignments_by_pk}
                    id={id}
                />
            </div>
        </div >
    )
}

const languages = (langs) => {
    const str = langs.reduce((s, v) => {
        return s + v.language.name + ", "
    }, "");

    return str.substring(0, str.length - 2);
}

const formatDate = (date) => {
    let d = new Date(date),
        month = '' + (d.getMonth() + 1),
        day = '' + d.getDate(),
        year = d.getFullYear();

    if (month.length < 2)
        month = '0' + month;
    if (day.length < 2)
        day = '0' + day;

    return [year, month, day].join('-');
}

export default AssignmentView;