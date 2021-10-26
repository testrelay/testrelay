import { useMutation, useQuery } from "@apollo/client";
import React, { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Loading } from "../../../components";
import { AlertError } from "../../../components/alerts";
import { UpdateReviewer } from "../../components/assignments/add-reviewer";
import { DELETE_REVIEWER, GET_ASSIGNMENT, INSERT_REVIEWER } from "../../components/assignments/queries";
import AssignmentStatus from "../../components/assignments/status";
import Timeline from "../../components/assignments/timeline";
import { dateToReadable, testLimitToReadable } from "../../../components/date";

const AssignmentView = () => {
	let { id } = useParams();
	const { loading, error, data } = useQuery(GET_ASSIGNMENT, { fetchPolicy: "network-only", variables: { id } });

	if (loading) {
		return <Loading />
	}

	if (error) {
		console.log("error fetching assignment", error)
		return <AlertError message="cant fetch assignment information, please reload the page" />
	}

	const scheduledFor = (assignment) => {
		if (assignment.test_time_chosen == null) {
			return (<span className="italic text-gray-400">unscheduled</span>);
		}

		return dateToReadable(new Date(assignment.test_day_chosen)) + " at " + assignment.test_time_chosen + " (" + assignment.test_timezone_chosen + ")";
	}

	const repoUrl = (url) => {
		if (url == null) {
			return (<span className="italic text-gray-400">repository not created, waiting for candidate to schedule</span>)
		}

		return (<a className="text-indigo-500" href={"https://" + url} rel="noreferrer" target="_blank">{"https://" + url}</a>)
	}

	return (
		<div>
			<div className="card shadow-md rounded-xl bg-white">
				<div className="p-4 pb-4 sm:px-8">
					<h2 className="text-xl text-indigo-500 capitalize mb-1">{data.assignments_by_pk.candidate_name}</h2>
					<h3 className="text-md text-gray-300 flex">
						<svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
						</svg>
						<span className="ml-2">{data.assignments_by_pk.candidate_email}</span>
					</h3>
				</div>
				<div class="border-t border-gray-200">
					<dl>
						<div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Test
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								<Link className="text-indigo-500" to={"/tests/" + data.assignments_by_pk.test.id + "/view"}>{data.assignments_by_pk.test.name}</Link>
							</dd>
						</div>
						<div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Assignment instructions
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								<a className="text-indigo-500" href={"https://github.com/" + data.assignments_by_pk.test.github_repo} rel="noreferrer" target="_blank">{"https://github.com/" + data.assignments_by_pk.test.github_repo}</a>
							</dd>
						</div>
						<div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Time limit
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								{testLimitToReadable(data.assignments_by_pk.time_limit)}
							</dd>
						</div>
						<div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Candidate to take test
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								{scheduledFor(data.assignments_by_pk)}
							</dd>
						</div>
						<div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Candidate Code
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								{repoUrl(data.assignments_by_pk.github_repo_url)}
							</dd>
						</div>
						<div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
							<dt class="text-sm font-medium text-gray-500">
								Current status
							</dt>
							<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
								<AssignmentStatus status={data.assignments_by_pk.status} />
							</dd>
						</div>
					</dl>
				</div>
			</div>
			<div>
				<div className="bg-white shadow-md rounded-xl mt-10">
					<div className>
						<h2 className="py-4 px-8 text-xl text-indigo-500 capitalize mb-1">Reviewers</h2>
					</div>
					<Reviewers reviewers={data.assignments_by_pk.reviewers} assignment_id={id} />
				</div>
			</div>

			<div className="flex flex-row">
				<div className="w-2/3">
					<Timeline id={id} events={data.assignments_by_pk.assignment_events} test={data.assignments_by_pk.test} />
				</div>
			</div>
		</div>
	);
}

const Reviewers = ({ reviewers, assignment_id }) => {
	const [users, setUsers] = useState([]);
	const [insertUser] = useMutation(INSERT_REVIEWER);
	const [deleteUser] = useMutation(DELETE_REVIEWER);

	const removeUser = (id) => {
		const filtered = users.reduce((arr, u) => {
			if (u.user.id !== id) {
				arr.push(u);
			}

			return arr;
		}, []);

		setUsers(filtered);
		deleteUser({ variables: { user_id: id } });
	}

	const addUser = (option) => {
		setUsers([...users, option]);
		insertUser({ variables: { user_id: option.user.id, assignment_id } });
	}

	useEffect(() => {
		const merged = arrayUnique([...users, ...reviewers]);
		setUsers(merged)
	}, [reviewers, users]);

	const formatted = users.map((r) => {
		let g = (
			<span className="text-md text-red-300 flex items-center">
				<svg className="h-4 w-4 text-red-300" width="20px" height="20px" viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
					<g>
						<path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" ></path>
					</g>
				</svg>
				<span className="ml-2">github account is not yet connected.</span>
			</span>
		)

		if (r.user.github_username) {
			g = (
				<span className="text-md text-gray-300 flex items-center">
					<svg className="h-4 w-4 text-gray-300" width="20px" height="20px" viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
						<g>
							<path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" ></path>
						</g>
					</svg>
					<span className="ml-2">{r.user.github_username}</span>)
				</span>
			)
		}

		return (
			<div key={r.user.id}>
				<div className="flex">
					<div className="flex-grow">
						<span className="text-gray-600 font-bold">{r.user.email}</span>
						{g}
					</div>
					<div>
						<svg xmlns="http://www.w3.org/2000/svg" onClick={() => { removeUser(r.user.id) }} className="h-6 w-6 text-gray-600 hover:text-indigo-500 cursor-pointer" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
						</svg>
					</div>
				</div>
			</div>
		)
	});

	return (<div className="border-t border-gray-200">
		<div className="space-y-3 py-4 px-8">
			{formatted}
		</div>
		<div className="bg-gray-50 py-4 px-8 rounded-b-lg">
			<UpdateReviewer selectedUsers={users} addUser={addUser} />
		</div>
	</div>);
}


function arrayUnique(array) {
	var a = array.concat();
	for (var i = 0; i < a.length; ++i) {
		for (var j = i + 1; j < a.length; ++j) {
			if (a[i].user.id === a[j].user.id)
				a.splice(j--, 1);
		}
	}

	return a;
}

export default AssignmentView;
