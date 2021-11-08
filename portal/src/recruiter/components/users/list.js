import { useQuery } from "@apollo/client";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Loading } from "../../../components";
import { useBusiness } from "../business/hook";
import EmptyState from "../../../components/empty-state";
import { GET_USERS } from "./queries";

const List = (props) => {
    const { selected } = useBusiness();
    const { data, loading: queryLoading } = useQuery(GET_USERS, {
        variables: {
            business_id: selected.id
        },
    });
    const [loading, setLoading] = useState(queryLoading);
    const [users, setUsers] = useState([]);

    useEffect(() => {
        setLoading(queryLoading);
    }, [queryLoading])

    useEffect(() => {
        if (data) {
            setUsers(data.users);
        }
    }, [data])

    if (loading) {
        return <Loading />
    }

    if (users.length === 1) {
        return (
            <EmptyState
                link="/users/create"
                icon="user"
                title="Invite your first user"
                description="Invite other users from your organisation to TestRelay. This is essential if you want others to review your candidate's completed assesments."
            />
        )
    }

    const content = users.map((e) => {
        return (
            <div key={e.id} className="p-2 lg:w-1/2 md:w-1/2 w-full">
                <div className="bg-white px-8 py-6 shadow-md rounded-lg">
                    <div className="flex flex-row">
                        <div className="w-2/3">
                            <h2 className="text-gray-900 title-font font-medium">{e.email}</h2>
                            <p className="text-gray-500">{e.business_users[0].user_type}</p>
                        </div>
                        <div className="w-1/3 flex items-center justify-end">
                            <Link to={"/users/" + e.id + "/view"}><svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                            </svg></Link>
                        </div>
                    </div>
                </div>
            </div>
        )
    })

    return (
        <div className="flex flex-wrap -m-2">
            {content}
        </div>
    )
}

export default List;