import React, { useEffect, useState } from "react";
import { getAuth, onIdTokenChanged } from "firebase/auth";
import firebase from "./firebase";
import react from "react";

const getClaims = async (user) => {
    const idTokenResult = await user.getIdTokenResult();
    return idTokenResult.claims["https://hasura.io/jwt/claims"];
}

const UserContext = react.createContext({ user: null, loading: true, token: "", claims: null });

const useUser = () => {
    const [state, setState] = useState({
        user: null,
        loading: true,
        token: "",
        claims: null
    });

    const auth = getAuth(firebase);

    useEffect(() => {
        const listen = onIdTokenChanged(auth, async (user) => {
            if (user) {
                const claims = await getClaims(user);
                const token = await user.getIdToken();
                setState(s => { return { ...s, loading: false, user, claims, token } });
                return;
            }

            setState(s => { return { ...s, loading: false, user, token: "", claims: null } });
        });

        return () => {
            listen();
        }
    }, [auth]);

    return state;
}

const FirebaseAuthProvider = ({ children }) => {
    const state = useUser();

    return (
        <UserContext.Provider value={state}>
            {children}
        </UserContext.Provider>
    );
}

const useFirebaseAuth = () => {
    const context = React.useContext(UserContext);

    if (context === undefined) {
        throw new Error(
            "useFirebaseAuth must be used within a FirebaseAuthProvider"
        );
    }

    return context;
}

const useFirebasetoken = () => {
    const context = React.useContext(UserContext);

    if (context === undefined) {
        throw new Error(
            "useFirebaseAuth must be used within a FirebaseAuthProvider"
        );
    }


    return async () => {
        const user = context.user;
        if (!user) {
            return context;
        }

        // if no claims refresh the token
        const claims = await getClaims(user);
        if (!claims) {
            return user.getIdToken(true);
        }

        return user.getIdToken();
    }
}

export { FirebaseAuthProvider, useFirebaseAuth, getClaims, useFirebasetoken };