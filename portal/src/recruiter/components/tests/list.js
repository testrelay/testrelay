import React, { useState } from "react";
import { Link } from 'react-router-dom';
import { GET_TESTS } from './queries';
import { useQuery } from "@apollo/client";
import { Loading } from "../../../components";
import { testLimitToReadable } from "../../../components/date";
import Pagination from "../../../components/pagination";
import EmptyState from "../../../components/empty-state";


const TestRow = (props) => {
    const code = props.zip || props.github_repo;
    return (
        <div key={props.id} className="bg-white relative shadow-md px-4 md:px-8 py-4 md:py-6 rounded text-center md:text-left mb-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="pb-4 md:pb-0 border-b-2 md:border-b-0">
                    <div className="capitalize text-md md:text-sm font-medium text-indigo-500 mb-2">
                        {props.name}
                    </div>
                    <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                        </svg>
                        <span className="ml-2">{code}</span>
                    </div>
                </div>
                <div className="pb-2 md:pb-0">
                    <div className="text-md md:text-sm font-medium text-gray-800 mb-2">
                        Default Time limit
                    </div>
                    <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span className="ml-2">{testLimitToReadable(props.time_limit)}</span>
                    </div>
                </div>
                <div className="flex items-center justify-center md:justify-end bg-indigo-500 md:bg-transparent text-white md:text-gray-800 p-2 md:p-0 rounded">
                    <Link to={"/tests/" + props.id + "/view"}><svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg></Link>
                </div>
            </div>
        </div>
    )
}

const Rows = () => {
    const [page, setPage] = useState(0);
    const limit = 2;

    const { loading, error, data } = useQuery(GET_TESTS, {
        fetchPolicy: "network-only",
        variables: {
            limit,
            offset: limit * page
        }
    })

    if (error) {
        console.log(error)
        return (null)
    }

    if (loading) {
        return <Loading />
    }


    if (data) {
        if (data.tests.length === 0) {
            return (
                <EmptyState
                    link="/tests/create"
                    icon="test"
                    title="Create your first test"
                    description="Tests represent a set of instructions that you'll send to a candidate. You'll need to create at least one before you can schedule an assignment for a candidate."
                />
            )
        }

        const rows = data.tests.map((e) => {
            return (<TestRow key={e.id} {...e} />)
        })

        return (
            <div>
                {rows}
                <Pagination setState={setPage} page={page} limit={limit} total={data.tests_aggregate.aggregate.count} />
            </div>
        )
    }

    return (null)
}
const TestList = () => {
    return (
        <Rows />
    )
};

export default TestList;