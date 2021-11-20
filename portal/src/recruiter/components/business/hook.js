import React, {useEffect, useState} from 'react';
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
    const [loading, setLoading] = useState(!stored);
    const {user, claims, loading: userLoading} = useFirebaseAuth(null);

    const id = claims ? parseInt(claims["x-hasura-user-pk"]) : null;
    const {data, error} = useQuery(GET_BUSINESS, {
        skip: !claims || selected,
        fetchPolicy: 'network-only',
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
            localStorage.removeItem('business');
        }

        setSelected(val);
    }


    useEffect(() => {
        if (error) {
            choose(null);
        }
    }, [error]);

    useEffect(() => {
        if (user == null && userLoading === false) {
            persistSelected(null);
        }
    }, [user, userLoading]);

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
        // eslint-disable
    }, [data, id])

    /* eslint-enable */

    return {loading, selected, setSelected: persistSelected};
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