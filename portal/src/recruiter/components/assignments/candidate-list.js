import {useQuery} from "@apollo/client";
import React, {useState} from "react";
import {Link} from 'react-router-dom';
import Loading from "../../../components/loading";
import {GET_ASSIGNMENTS} from "./queries";
import AssignmentStatus from "./status";
import Pagination from "../../../components/pagination";
import EmptyState from "../../../components/empty-state";
import {useBusiness} from "../business/hook";

const CandidateRow = (props) => {
    return (
        <div key={props.id}
             className="bg-white relative shadow shadow-md px-8 py-6 rounded text-center md:text-left mb-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="pb-4 md:pb-0 border-b-2 md:border-b-0">
                    <div className="text-md md:text-sm font-medium text-indigo-500 mb-2">
                        {props.name}
                    </div>
                    <div className="text-sm text-gray-500 flex items-center justify-center md:justify-start">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24"
                             stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                        </svg>
                        <span className="ml-1">{props.email}</span>
                    </div>
                </div>
                <div className="pb-4 md:pb-0">
                    <div className="text-md md:text-sm font-medium text-gray-800 mb-2">
                        {props.test}
                    </div>
                    <div>
                        <AssignmentStatus status={props.status}/>
                    </div>
                </div>
                <div
                    className="flex items-center justify-center md:justify-end bg-indigo-500 md:bg-transparent text-white md:text-gray-800 p-2 md:p-0 rounded">
                    <Link to={"/assignments/" + props.id + "/view"}>
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24"
                             stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7"/>
                        </svg>
                    </Link>
                </div>
            </div>
        </div>
    )
}

const Rows = () => {
    const [page, setPage] = useState(0);
    const {selected} = useBusiness();
    const limit = 2;
    const {loading, error, data} = useQuery(GET_ASSIGNMENTS, {
        fetchPolicy: "network-only",
        variables: {
            limit,
            offset: limit * page,
            business_id: selected.id
        }
    })

    if (loading) {
        return <Loading/>
    }

    if (error) {
        console.log(error)
        return (null)
    }

    if (data) {
        if (data.assignments.length === 0) {
            return (
                <EmptyState
                    link="/assignments/create"
                    icon="assignment"
                    title="Schedule your first assignment"
                    description="Assignments are the way you can invite a candidate to take one of your company's tests. You'll be able to track a candidates progress & review the completed code."
                />)
        }


        const rows = data.assignments.map((e) => {
            return (
                <CandidateRow key={e.id} id={e.id} name={e.candidate_name} email={e.candidate_email} status={e.status}
                              test={e.test.name}/>)
        })

        return (
            <div>
                {rows}
                <Pagination setState={setPage} page={page} limit={limit}
                            total={data.assignments_aggregate.aggregate.count}/>
            </div>
        )
    }

    return (null)
}

const CandidateList = () => {
    return (
        <Rows/>
    )
};

export default CandidateList;