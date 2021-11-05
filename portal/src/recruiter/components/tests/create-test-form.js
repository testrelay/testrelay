import React, { useEffect, useState } from "react";
import { Redirect } from 'react-router-dom';
import { useQuery, useMutation } from '@apollo/client';
import { GET_REPOS, INSERT_REPO } from './queries';
import { Loading } from '../../../components';
import CodeSelect from "./code-select";
import { useBusiness } from "../business/hook";

const GithubSelect = (props) => {
	const hasAuthed = props.business.github_installation_id != null;
	const { reposLoading, error, data } = useQuery(GET_REPOS, { skip: !hasAuthed, variables: { id: props.business.id } })

	if (!hasAuthed) {
		return (
			<div>
				<a href="https://github.com/apps/testrelay" className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto inline-flex items-center">
					Connect github
					<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="inline-block w-6 h-6 ml-2 stroke-current">
						<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7"></path>
					</svg>
				</a>
			</div>
		)
	}

	if (reposLoading) {
		return (<Loading />)
	}

	if (error) {
		return (
			<div className="alert alert-error">
				<div className="flex-1">
					<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="w-6 h-6 mx-2 stroke-current">
						<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
					</svg>
					<label>We could not fetch your github repos, please check the app is installed correctly</label>
				</div>
			</div>
		)
	}

	if (data) {
		const options = data.repos.map((e, i) => {
			return (<option key={i} value={e.full_name}>{e.full_name}</option>)
		});

		return (
			<select name="githubRepo" className="select select-bordered block appearance-none w-full" onChange={props.change}>
				<option value="" disabled selected>Select your repo</option>
				{options}
			</select>
		)
	}

	return (<Loading />)
}

const CreateTestForm = (props) => {
	const { selected } = useBusiness();
	const [testForm, setTestForm] = useState({
		name: null,
		timeLimit: 1,
		timeLimitUnit: 'hours',
		testWindow: 1,
		testWindowUnit: 'days',
		githubRepo: null,
		languages: []
	});

	const [insertRepo, mutation] = useMutation(INSERT_REPO)
	const insertLoading = mutation.loading;
	const insertData = mutation.data;
	const insertError = mutation.error;

	const [error, setError] = useState(null);

	useEffect(() => {
		if (insertError) {
			setError(insertError);
		}
	}, [insertError]);

	const optionSelection = () => {
		if (testForm['timeLimitUnit'] === "days") {
			return 30;
		}

		return 24;
	}

	const inputChange = (event) => {
		setTestForm(Object.assign({}, testForm, { [event.target.name]: event.target.value }));
	}

	const setLanguages = (languages) => {
		setTestForm(Object.assign({}, testForm, { languages }))
	}

	const submitForm = (event) => {
		event.preventDefault();

		const windowsInSecs = {
			'hours': 3600,
			'days': 86400,
			'weeks': 604800
		};

		const testExpires = parseInt(testForm.testWindow) * windowsInSecs[testForm.testWindowUnit];

		insertRepo({
			variables: {
				github_repo: testForm.githubRepo,
				name: testForm.name,
				test_window: testExpires,
				time_limit: parseInt(testForm.timeLimit) * windowsInSecs[testForm.timeLimitUnit],
				languages: testForm.languages,
				business_id: selected.id
			}
		}).catch(e => {
			console.log("insert test error", e);
			setError(e);
		})
	}

	if (insertData) {
		return (<Redirect push to="/tests/" />)
	}

	return (
		<div className="pb-8">
			<div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
				<label className="block text-gray-700 text-md font-bold mb-2">
					Test Name
				</label>
				<input onChange={inputChange} name="name" className="input input-bordered w-full md:w-1/2 py-2 px-3 text-gray-700" type="text" placeholder="e.g. BE candidate Test" />
			</div>
			<div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
				<label className="block text-gray-700 text-md font-bold mb-2">Upload the candidate test instructions/code</label>
					<div className="w-full">
						<p className="mb-4 text-sm">Select a repository to use as a test template.
							TestRelay will clone this repository and upload the files to the candidates test on commencement. </p>
						<GithubSelect business={selected} change={inputChange} />
					</div>
			</div>
			<div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
				<label className="block tracking-wide text-gray-700 text-md font-bold mb-1">
					Test Languages
				</label>
				<p className="mb-3 text-sm">Select one or more languages that candidates can complete the test in. (Select or type to create options).</p>
				<CodeSelect setState={setLanguages} />
			</div>
			<div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
				<div className="flex flex-wrap -mx-3 mb-8 items-baseline">
					<div className="w-full px-3 mb-3">
						<label className="block tracking-wide text-gray-700 text-md font-bold mb-1">
							Candidate time limit
						</label>
						<p className="text-sm">The amount of time a candidate has to complete the test. Starting from when TestRelay uploads the test code to the candidate repository. This is a default value and can be changed later for individual candidates.</p>
					</div>
					<div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
						<div className="relative">
							<select name="timeLimit" onChange={inputChange} className="select select-bordered block appearance-none w-full">
								{optionsTo(optionSelection())}
							</select>
						</div>
					</div>
					<div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
						<div className="relative">
							<select name="timeLimitUnit" onChange={inputChange} className="select select-bordered block appearance-none w-full">
								<option value="hours">hours</option>
								<option value="days">days</option>
							</select>
						</div>
					</div>
				</div>
				<div className="flex flex-wrap -mx-3">
					<div className="w-full px-3 mb-3">
						<label className="block tracking-wide text-gray-700 font-bold text-md mb-1">
							Test Expiry
						</label>
						<p className="text-sm">The time from sending a test invite email that a candidate has to schedule a time to sit the test. This is a default value and can be changed later for individual candidates.</p>
					</div>
					<div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
						<div className="relative">
							<select name="testWindow" onChange={inputChange} className="select select-bordered block appearance-none w-full">
								{optionsTo(30)}
							</select>
						</div>
					</div>
					<div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
						<div className="relative">
							<select name="testWindowUnit" onChange={inputChange} className="select select-bordered block appearance-none w-full">
								<option value="days">days</option>
								<option value="weeks">weeks</option>
							</select>
						</div>
					</div>
				</div>
			</div>
			<div className="w-full bg-white px-8 py-6 mb-8 shadow-md rounded">
				<SubmitBtn loading={insertLoading} submit={submitForm} />
				{error &&
					<div class="alert alert-error">
						<div class="flex-1">
							<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-6 h-6 mx-2 stroke-current">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
							</svg>
							<label>Could not create test, please try again</label>
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
		<button className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto" onClick={props.submit}>
			Save Test
		</button>
	)
}

const optionsTo = (num) => {
	let options = [];

	for (let i = 1; i <= num; i++) {
		options.push(<option value={i} key={i}>{i}</option>)
	}

	return options;
}

export default CreateTestForm;