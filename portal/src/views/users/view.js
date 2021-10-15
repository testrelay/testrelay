import React from "react";
import { useQuery } from "@apollo/client";
import { useParams } from "react-router";
import { Loading } from "../../components";
import { GET_USER } from "../../components/users/queries";

const View = () => {
    const { id } = useParams();
    const { loading, data } = useQuery(GET_USER, {
        variables: {
            id
        }
    });

    if (loading) {
        return (<Loading />);
    }


    const username = (data.users_by_pk.github_username) ? (<span className="text-indigo-500">{data.users_by_pk.github_username}</span>) : (<span className="text-indigo-500">github account not connected</span>);

    return (
        <div className="card  shadow-md rounded-xl bg-white">
            <div className="p-4 pb-4 sm:px-8">
                <h3 className="text-md text-gray-500 flex">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                    <span className="ml-2">{data.users_by_pk.email}</span>
                </h3>
            </div>
            <div class="border-t border-gray-200">
                <dl>
                    <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-8">
                        <dt class="text-sm font-medium text-gray-500">
                            Github Username
                        </dt>
                        <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                            {username}
                        </dd>
                    </div>
                </dl>
            </div>
        </div>
    )
}

export default View;