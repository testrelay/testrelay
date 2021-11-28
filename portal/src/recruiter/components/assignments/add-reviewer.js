import React, {useEffect, useState} from 'react';
import CreatableSelect from 'react-select/creatable';
import {GET_USERS, INVITE_USER} from "../users/queries";
import {useBusiness} from "../business/hook";
import {useMutation, useQuery} from '@apollo/client';
import {AlertError} from '../../../components/alerts';

const AddReviewer = (props) => {
    const [loading, setLoading] = useState(false);
    const [users, setUsers] = useState([]);
    const [error, setError] = useState(null);
    const [vals, setVals] = useState([]);
    const {selected} = useBusiness();
    const [inviteUser, {data: muData, loading: muLoading, error: muError}] = useMutation(INVITE_USER);

    const {data, loading: queryLoading} = useQuery(GET_USERS, {
        fetchPolicy: 'network-only',
        variables: {
            business_id: selected.id
        },
    });

    useEffect(() => {
        if (data) {
            const options = data.users.map((e) => {
                return {value: e.id, label: e.email};
            });

            setUsers(options);
        }
    }, [data]);

    useEffect(() => {
        if (muData) {
            setVals(v => [...v, {id: muData.inviteUser.id}]);
            setUsers(v => [...v, {value: muData.inviteUser.id, label: muData.inviteUser.email}]);
        }
    }, [muData]);

    useEffect(() => {
        if (muError) {
            setError("could not invite user please refresh and try again");
            setLoading(false);
        }
    }, [muError])

    useEffect(() => {
        setLoading(queryLoading || muLoading);
    }, [queryLoading, muLoading]);

    const handleChange = (options) => {
        setVals(options);

        props.reviewerChange(options.map(e => {
            return {user_id: e.value}
        }));
    };

    const handleCreate = async (value) => {
        try {
            await inviteUser({
                variables: {
                    business_id: selected.id,
                    email: value,
                    redirect_link: process.env.REACT_APP_URL + "/assignments/assigned"
                }
            });
        } catch (e) {

        }
    }


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
            <AlertError message={error}/>
            }
        </div>
    );
}


const UpdateReviewer = ({selectedUsers, addUser}) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const {selected} = useBusiness();
    const [inviteUser, {data: muData, loading: muLoading, error: muError}] = useMutation(INVITE_USER);
    const [options, setOptions] = useState([]);
    const [lookup, setLookup] = useState({});

    const {data, loading: queryLoading} = useQuery(GET_USERS, {
        fetchPolicy: 'network-only',
        variables: {
            business_id: selected.id
        },
    });

    useEffect(() => {
        const lookup = selectedUsers.reduce((ob, u) => {
            ob[u.user.id] = u;

            return ob;
        }, {});

        setLookup(lookup);
    }, [selectedUsers])

    useEffect(() => {
        if (data != null) {
            const options = data.users.reduce((arr, e) => {
                if (lookup[e.id] == null) {
                    arr.push({value: e.id, label: e.email});
                }

                return arr;
            }, []);

            setOptions(options);
        }
    }, [data, lookup])

    useEffect(() => {
        if (muData) {
            addUser({user: {id: muData.inviteUser.id, email: muData.inviteUser.email, github_username: null}})
        }
    }, [muData, addUser]);

    useEffect(() => {
        if (muError) {
            setError("could not invite user please refresh and try again");
            setLoading(false);
        }
    }, [muError])

    useEffect(() => {
        setLoading(queryLoading || muLoading);
    }, [queryLoading, muLoading]);

    const handleChange = (option) => {
        const user = data.users.find((ob) => {
            return ob.id === option.value;
        }, {});

        addUser({user});
    };

    const handleCreate = async (value) => {
        try {
            const exists = options.find(e => {
                return e.label === value;
            });

            if (exists != null) {
                return;
            }

            await inviteUser({
                variables: {
                    business_id: selected.id,
                    email: value,
                    redirect_link: process.env.REACT_APP_URL + "/assignments/assigned"
                }
            });
        } catch (e) {
        }
    }

    return (
        <div className="mb-2">
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
            <AlertError message={error}/>
            }
        </div>
    );
}


export {AddReviewer, UpdateReviewer};