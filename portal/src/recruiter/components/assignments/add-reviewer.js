import React, { useEffect, useState } from 'react';
import CreatableSelect from 'react-select/creatable';
import { GET_USERS } from "../users/queries";
import { getFunctions, httpsCallable } from "@firebase/functions";
import firebase from "../../../auth/firebase";
import { useBusiness } from "../business/hook";
import { useQuery } from '@apollo/client';
import { AlertError } from '../../../components/alerts';

const AddReviewer = (props) => {
    const [loading, setLoading] = useState(false);
    const [users, setUsers] = useState([]);
    const [error, setError] = useState(null);
    const [vals, setVals] = useState([]);
    const { selected } = useBusiness();

    const { data, loading: queryLoading } = useQuery(GET_USERS, {
        fetchPolicy: 'network-only',
        variables: {
            business_id: selected.id
        },
    });

    useEffect(() => {
        if (data) {
            const options = data.users.map((e) => {
                return { value: e.id, label: e.email };
            })

            setUsers(options);
        }
    }, [data]);

    useEffect(() => {
        setLoading(queryLoading);
    }, [queryLoading]);

    const handleChange = (options) => {
        setVals(options);

        props.reviewerChange(options.map(e => { return { user_id: e.value } }));
    };

    const handleCreate = async (value) => {
        setLoading(true);

        const functions = getFunctions(firebase, "europe-west2");
        const invite = httpsCallable(functions, "inviteUser");

        try {
            const data = await invite({ email: value, business_name: selected.name, business_id: selected.id });
            console.log('data', data);
            setVals([...vals, { id: data.id }])
            setUsers([...vals, { value: data.id, label: value }])
            setLoading(false);
        } catch (error) {
            setError("could not invite user " + value + " please refresh and try again");
            setLoading(false);
        }
    };


    return (
        <div>
            <CreatableSelect
                placeholder="Add reviewer"
                isMulti
                isDisabled={loading}
                isLoading={loading}
                onChange={handleChange}
                onCreateOption={handleCreate}
                options={users}
                value={vals}
                styles={{
                    control: (provided) => {
                        return {
                            ...provided,
                            padding: "0.4rem 0.3rem",
                        }
                    }
                }}
            />
            {error &&
                <AlertError message={error} />
            }
        </div>
    );
}


const UpdateReviewer = ({ selectedUsers, addUser }) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [options, setOptions] = useState([]);
    const { selected } = useBusiness();

    const { data, loading: queryLoading } = useQuery(GET_USERS, {
        fetchPolicy: 'network-only',
        variables: {
            business_id: selected.id
        },
    });

    const lookup = selectedUsers.reduce((ob, u) => {
        ob[u.user.id] = u;

        return ob;
    }, {});

    useEffect(() => {
        if (data) {
            const options = data.users.reduce((arr, e) => {
                if (lookup[e.id] == null) {
                    arr.push({ value: e.id, label: e.email });
                }

                return arr;
            }, [])

            console.log(options);
            setOptions(options);
        }
    }, [selectedUsers, data, lookup]);

    useEffect(() => {
        setLoading(queryLoading);
    }, [queryLoading]);

    const handleChange = (option) => {
        const user = data.users.find((ob) => {
            return ob.id === option.value;
        }, {});

        addUser({ user });
    };

    const handleCreate = async (value) => {
        setLoading(true);

        const functions = getFunctions(firebase, "europe-west2");
        const invite = httpsCallable(functions, "inviteUser");

        try {
            const data = await invite({ email: value, business_name: selected.name, business_id: selected.id });
            addUser({ user: { id: data.id, email: value, github_username: null } })
            setLoading(false);
        } catch (error) {
            console.log(error);
            setError("could not invite user " + value + " please refresh and try again");
            setLoading(false);
        }
    };

    return (
        <div>
            <CreatableSelect
                placeholder="Add reviewer"
                isDisabled={loading}
                isLoading={loading}
                onChange={handleChange}
                onCreateOption={handleCreate}
                options={options}
                value=""
                styles={{
                    control: (provided) => {
                        return {
                            ...provided,
                            padding: "0.4rem 0.3rem",
                        }
                    }
                }}
            />
            {error &&
                <AlertError message={error} />
            }
        </div>
    );
}



export { AddReviewer, UpdateReviewer };