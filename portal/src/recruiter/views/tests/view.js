import { useQuery } from "@apollo/client";
import React from "react";
import { useParams, Link } from "react-router-dom";
import { Loading } from "../../../components";
import { AlertError } from "../../../components/alerts";
import { testLimitToReadable } from "../../../components/date";
import EmptyState from "../../../components/empty-state";
import Languages from "../../components/tests/languages";
import { GET_TEST } from "../../components/tests/queries";

const Assignments = (props) => {
    if (props.assignments.length === 0) {
        return (
            <EmptyState
                link="/assignments/create"
                icon="assignment"
                title="Schedule an assignment"
                description="You have no assignments using this test. Create one now and invite a candidate to interview."
            />)
    }
    const asg = props.assignments.map((e, i) => {
        return (
            <div key={i} className="bg-white relative shadow shadow-md px-8 py-4 rounded text-center md:text-left">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <div className="text-md md:text-sm font-medium text-indigo-500 mb-1">
                            {e.candidate_name}
                        </div>
                        <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                            </svg>
                            <span className="ml-1">{e.candidate_email}</span>
                        </div>
                    </div>
                    <div className="flex items-center justify-end">
                        <Link to={"/assignments/" + e.id + "/view"}><svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg></Link>
                    </div>
                </div>
            </div>
        )
    })

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {asg}
        </div>
    )
}

const TestView = () => {
    let { id } = useParams();
    const { loading, error, data } = useQuery(GET_TEST, { fetchPolicy: "network-only", variables: { id } });

    if (loading) {
        return <Loading />
    }

    if (error) {
        return <AlertError message="cant fetch test information, please reload the page" />
    }

    return (
        <div>
            <div className="shadow-md rounded mb-6 bg-white">
                <div className="p-4 pb-4 sm:px-8">
                    <h2 className="text-xl text-indigo-500 capitalize mb-1">{data.tests_by_pk.name}</h2>
                    <h3 className="text-md text-gray-300 flex">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                        </svg>
                        <span className="ml-2">{data.tests_by_pk.github_repo}</span>
                    </h3>
                </div>
                <div className="border-t border-gray-200">
                    <dl>
                        <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
                            <dt className="text-sm font-medium text-gray-500">
                                Default time limit
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                                {testLimitToReadable(data.tests_by_pk.time_limit)}
                            </dd>
                        </div>
                        <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
                            <dt className="text-sm font-medium text-gray-500">
                                Default test expiry
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                                {testLimitToReadable(data.tests_by_pk.test_window)}
                            </dd>
                        </div>
                        <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
                            <dt className="text-sm font-medium text-gray-500 mb-2 sm:mb-0">
                                Allowed languages
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                                <Languages languages={data.tests_by_pk.test_languages} />
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>
            <div className="py-4 border-b-4 mb-6">
                <h2 className="text-xl font-bold">Test Assignments</h2>
            </div>
            <Assignments assignments={data.tests_by_pk.assignments} />
        </div>
    );
}

export default TestView;
