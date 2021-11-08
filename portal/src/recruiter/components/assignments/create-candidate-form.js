import React, {useState} from "react";
import {Redirect} from 'react-router-dom';
import {useMutation, useQuery} from '@apollo/client';
import {GET_TESTS} from '../tests/queries';
import {INSERT_ASSIGNMENT} from './queries.js';
import {Loading} from '../../../components';
import {dateToReadable, formatDate} from "../../../components/date";
import DatePicker from 'react-datepicker';
import {AddReviewer} from "./add-reviewer";

const CreateCandidateForm = (props) => {
    const today = new Date();
    const [candidateForm, setCandidateForm] = useState({
        name: null,
        emai: null,
        test: null,
        reviewers: [],
        timeLimit: null,
        timeLimitUnit: 'hours',
        testWindowUnit: 'days',
        alterTest: false,
        chooseUntil: today.setDate(today.getDate() + 1)
    });


    const {loading, data} = useQuery(GET_TESTS);

    const [insertCandidate, mutation] = useMutation(INSERT_ASSIGNMENT)
    const insertLoading = mutation.loading;
    const insertData = mutation.data;
    const insertError = mutation.error;

    const inputChange = (event) => {
        let value = event.target.value;
        if (event.target.type === 'checkbox') {
            value = event.target.checked;
        }

        setCandidateForm(Object.assign({}, candidateForm, {[event.target.name]: value}));
    }

    const setReviewers = (reviewers) => {
        setCandidateForm(Object.assign({}, candidateForm, {reviewers}));
    }

    const dateChange = (date) => {
        setCandidateForm(Object.assign({}, candidateForm, {chooseUntil: date}));
    }

    const testChange = (event) => {
        const test = data.tests.find((e) => {
            return e.id === parseInt(event.target.value);
        });

        setCandidateForm(Object.assign({}, candidateForm, {test}));
    }

    const submitForm = (event) => {
        event.preventDefault();

        const windowsInSecs = {
            'hours': 3600,
            'days': 86400,
            'months': 2629800
        };

        const test = candidateForm.test;

        let testWindow = test.test_window;
        let chooseUntil = new Date();
        chooseUntil = chooseUntil.setSeconds(chooseUntil.getSeconds() + testWindow)

        let timeLimit = test.time_limit;
        if (candidateForm.chooseUntil) {
            chooseUntil = candidateForm.chooseUntil;
        }

        if (candidateForm.timeLimit) {
            timeLimit = parseInt(candidateForm.timeLimit) * windowsInSecs[candidateForm.timeLimitUnit];
        }

        insertCandidate({
            variables: {
                name: candidateForm.name,
                email: candidateForm.email,
                test_id: candidateForm.test.id,
                choose_until: formatDate(chooseUntil),
                time_limit: timeLimit,
                reviewers: candidateForm.reviewers
            }
        }).catch(e => {
            console.log(e);
        })
    }

    if (insertData) {
        return (<Redirect push to={{
            pathname: "/assignments/",
            state: {success: "scheduling test for candidate " + insertData.insert_assignments_one.candidate_email}
        }}/>)
    }

    return (
        <div className="pb-12">
            <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                <div className="grid sm:grid-cols-2 gap-4">
                    <div className="">
                        <label className="block  text-gray-700 text-md font-bold mb-2">
                            Candidate Name
                        </label>
                        <input onChange={inputChange} name="name" className="input input-bordered rounded w-full"
                               type="text" placeholder="Joe Bloggs"/>
                    </div>
                    <div className="">
                        <label className="block  text-gray-700 text-md font-bold mb-2">
                            Candidate Email
                        </label>
                        <input onChange={inputChange} name="email" className="input input-bordered rounded w-full"
                               type="email" placeholder="joe@bloggs.com"/>
                    </div>
                </div>
            </div>
            <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                <label className="block  text-gray-700 text-md font-bold mb-2">
                    Candidate Test
                </label>
                <p className="mb-3 text-sm">Which test the candidate should sit.</p>
                <TestSelect loading={loading} data={data} change={testChange}/>
            </div>
            <TestDisplay test={candidateForm.test} change={inputChange}/>
            <AlterTestSection test={candidateForm.test} alter={candidateForm.alterTest} change={inputChange}
                              changeDate={dateChange} chooseUntil={candidateForm.chooseUntil}
                              timeLimitUnit={candidateForm.timeLimitUnit}/>
            <div className="w-full bg-white px-8 py-6 shadow-md rounded mb-8">
                <label className="block  text-gray-700 text-md font-bold mb-2">
                    Assignment Reviewers
                </label>
                <p className="mb-3 text-sm">Add users from your organisation that will review the code of this
                    assignment. These
                    users will need to have a github account to be able to access the assignment repository. You can add
                    more users later.</p>
                <AddReviewer reviewerChange={setReviewers}/>
            </div>
            <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
                <SubmitBtn loading={insertLoading} submit={submitForm} change={inputChange}/>
                {insertError &&
                <div className="alert alert-error mt-4">
                    <div className="flex-1">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                             className="w-6 h-6 mx-2 stroke-current">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                  d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
                        </svg>
                        <label>Could not create assignment, please try again</label>
                    </div>
                </div>
                }
            </div>
        </div>
    )
}


const SubmitBtn = (props) => {
    if (props.loading) {
        return (
            <button className="disabled bg-gray-600 text-white text-sm rounded px-4 py-2 w-auto">
                Loading
            </button>
        )
    }

    return (
        <button className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto"
                onClick={props.submit}>
            Schedule assignment
        </button>
    )
}

const TestDisplay = (props) => {
    if (props.test == null) {
        return null;
    }

    let unit = 'hours';
    let limit = Math.floor(props.test.time_limit / 3600);
    if (limit > 24) {
        unit = 'days';
        limit = limit / 24;
    }
    const d = new Date();
    d.setSeconds(d.getSeconds() + props.test.test_window)
    const window = dateToReadable(d)

    return (
        <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
            <p className="text-sm">The candidate will have <b>{limit} {unit}</b> to complete test hosted @ <a
                href={"https://github.com/" + props.test.github_repo} target="_blank" rel="noreferrer"
                className="text-indigo-500">github.com/{props.test.github_repo}</a></p>
            <p className="mb-4 text-sm">The will be able to schedule a time to take the test until <b>{window}</b>.</p>
            <div className="max-w-sm">
                <div className="form-control">
                    <label className="cursor-pointer flex items-center">
                        <input type="checkbox" name="alterTest" onChange={props.change}
                               className="checkbox checkbox-primary"/>
                        <span className="ml-2 label-text ">Change test default time constraints</span>
                    </label>
                </div>
            </div>
        </div>
    )
}

const AlterTestSection = (props) => {
    const optionsTo = (num) => {
        let options = [];

        for (let i = 1; i <= num; i++) {
            options.push(<option value={i} key={i}>{i}</option>)
        }

        return options;
    }

    const optionSelection = () => {
        if (props.timeLimitUnit === "days") {
            return 30;
        }

        return 24;
    }

    if (!props.alter) {
        return (null);
    }

    return (
        <div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
            <div className="grid gap-6 divide-y-2">
                <div>
                    <div className="mb-2">
                        <label className="block  tracking-wide text-gray-700 text-sm font-bold mb-1">
                            Candidate assignment time limit
                        </label>
                        <p className="text-sm">The amount of time a candidate has to complete the assignment.</p>
                    </div>
                    <div className="grid gap-4 sm:grid-cols-2">
                        <select name="timeLimit" onChange={props.change} className="select select-bordered w-full">
                            <option value="" disabled selected>Select time limit</option>
                            {optionsTo(optionSelection())}
                        </select>
                        <select name="timeLimitUnit" onChange={props.change}
                                className="select select-bordered block appearance-none w-full">
                            <option value="hours">hours</option>
                            <option value="days">days</option>
                        </select>
                    </div>
                </div>
                <div className="pt-6">
                    <label className="block  tracking-wide text-gray-700 font-bold text-sm mb-1">
                        Assignment Expiry
                    </label>
                    <p className="text-sm mb-2">The date that a candidate has until to take the assignment.</p>
                    <DatePicker
                        selected={props.chooseUntil}
                        onChange={(date) => props.changeDate(date)}
                        minDate={new Date()}
                        nextMonthButtonLabel=">"
                        previousMonthButtonLabel="<"
                    />
                </div>
            </div>
        </div>
    )
}

const TestSelect = (props) => {
    if (props.loading) {
        return <Loading/>
    }

    const tests = props.data.tests.map((e) => {
        return (<option key={e.id} value={e.id}>{e.name}</option>);
    })

    return (
        <select name="testId" onChange={props.change} className="select select-bordered w-full">
            <option value="" disabled selected>Select test</option>
            {tests}
        </select>
    )
}

export default CreateCandidateForm;