import React, {useCallback, useEffect, useState} from 'react';
import {useQuery} from '@apollo/client';
import {GET_BUSINESS} from './queries';
import {useFirebaseAuth} from '../../../auth/firebase-hooks';


const BusinessContext = React.createContext({
    selected: null, loading: false, setSelected: (val) => {
    }
});

const getSelected = () => {
    const item = localStorage.getItem('business');
    if (!item) {
        return null;
    }

    return JSON.parse(item);
}

const useBusinesses = () => {
    const stored = getSelected();

    const [selected, setSelected] = useState(stored);
    const [loading, setLoading] = useState(true);
    const {user, claims, loading: userLoading} = useFirebaseAuth(null);

    const id = claims ? parseInt(claims["x-hasura-user-pk"]) : null;
    const {data, error, loading: businessLoading} = useQuery(GET_BUSINESS, {
        skip: !claims || selected,
        fetchPolicy: 'network-only',
    });

    const choose = useCallback((val) => {
        if (val) {
            localStorage.setItem('business', JSON.stringify(val));
        } else {
            localStorage.removeItem('business');
        }

        setSelected(val);
    }, [setSelected]);

    useEffect(() => {
        if (error) {
            choose(null);
        }
    }, [choose, error]);

    useEffect(() => {
        if (selected) {
            setLoading(false);
        }
    }, [selected]);

    useEffect(() => {
        if (user == null && userLoading === false) {
            choose(null);
        }
    }, [choose, user, userLoading]);

    useEffect(() => {
        if (businessLoading) {
            setLoading(true);
        }
    }, [businessLoading]);

    useEffect(() => {
        if (data) {
            const returned = data.businesses || [];
            if (returned.length === 0) {
                choose(null);
            }

            const created = returned.find(e => e.creator_id === id)
            if (created) {
                choose(created);
            } else {
                choose(returned[0]);
            }
        }
    }, [choose, data, id])

    return {loading, selected, setSelected: choose};
}

const BusinessProvider = ({children}) => {
    const state = useBusinesses();

    return (
        <BusinessContext.Provider value={state}>
            {children}
        </BusinessContext.Provider>
    );
}

const useBusiness = () => {
    const context = React.useContext(BusinessContext);

    if (context === undefined) {
        throw new Error(
            "useBusiness must be used within a BusinessProvider"
        );
    }

    return context;
}

export {BusinessProvider, useBusiness}