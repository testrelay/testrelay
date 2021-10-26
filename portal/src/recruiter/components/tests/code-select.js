import React, { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client';

import CreatableSelect from 'react-select/creatable';
import { GET_LANGUAGES, INSERT_LANGUAGE } from './queries';

const CodeSelect = (props) => {
    const [insertLang] = useMutation(INSERT_LANGUAGE);

    const [state, localState] = useState({
        options: [], isLoading: true, value: []
    });

    const setState = (value) => {
        localState(value);

        const languages = value.value.map(e => { return { language_id: e.value } });
        props.setState(languages);
    }

    useQuery(GET_LANGUAGES, {
        fetchPolicy: "network-only",
        onCompleted: (data) => {
            const options = data.languages.map((e) => {
                return { value: e.id, label: e.name };
            })

            setState(Object.assign({}, state, { options, isLoading: false }));
        },
    });

    const handleChange = (options, actionMeta) => {
        setState(Object.assign({}, state, {
            value: options,
        }));
    };

    const handleCreate = (inputValue) => {
        setState(Object.assign({}, state, { isLoading: true }));

        insertLang({
            variables: {
                name: inputValue.toLowerCase().replace(/\W/g, ''),
            }
        }).then((data) => {
            const newOption = {
                value: data.data.insert_languages_one.id,
                label: data.data.insert_languages_one.name
            };

            setState(Object.assign({}, state, {
                isLoading: false,
                options: [...state.options, newOption],
                value: [...state.value, newOption],
            }));
        })
    };

    return (
        <CreatableSelect
            isMulti
            isDisabled={state.isLoading}
            isLoading={state.isLoading}
            onChange={handleChange}
            onCreateOption={handleCreate}
            options={state.options}
            value={state.value}
            styles={{
                control: (provided) => {
                    return {
                        ...provided,
                        padding: "0.4rem 0.3rem",
                    }
                }
            }}
        />
    );
}

export default CodeSelect;