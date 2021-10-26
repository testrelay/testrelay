import React from 'react';
import { useQuery } from '@apollo/client';
import { useState, useEffect } from 'react';
import { GET_BUSINESS } from './queries';
import { useFirebaseAuth } from '../../auth/firebase-hooks';


const BusinessContext = React.createContext({ selected: null, loading: false, setSelected: (val) => { } });

const getSelected = () => {
    const item = localStorage.getItem('business');
    if (!item) {
        return null;
    }

    return JSON.parse(item);
}

const useBusinesses = () => {
    const stored = getSelected();
    console.log("stored biz in local storage", stored);

    const [selected, setSelected] = useState(stored);
    const [loading, setLoading] = useState(!stored);
    const { claims } = useFirebaseAuth(null);

    const id = claims ? parseInt(claims["x-hasura-user-pk"]) : null;
    const { data, error } = useQuery(GET_BUSINESS, {
        skip: !claims || selected,
        nextFetchPolicy: 'network-only',
    });

    /* eslint-disable */
    const choose = (val) => {
        persistSelected(val);
        setLoading(false);
    }

    const persistSelected = (val) => {
        if (val) {
            localStorage.setItem('business', JSON.stringify(val));
        } else {
            localStorage.setItem('business', val);
        }
        setSelected(val);
    }


    useEffect(() => {
        if (error) {
            console.log("master business error", error)
            choose(null);
        }
    }, [error]);


    useEffect(() => {
        if (data) {
            console.log("businesses returned from business provider ", data)
            const returned = data.businesses || [];
            if (returned.length === 0) {
                return choose(null);
            }

            const created = returned.find(e => e.creator_id === id)
            if (created) {
                return choose(created);
            }

            choose(returned[0]);
        }
        // eslint-disable
    }, [data, id])

    /* eslint-enable */

    return { loading, selected, setSelected: persistSelected };
}

const BusinessProvider = ({ children }) => {
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

export { BusinessProvider, useBusiness }