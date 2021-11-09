import React from "react";
import react, {useEffect, useState} from "react";
import {getAuth, onIdTokenChanged} from "firebase/auth";
import {getFunctions, httpsCallable} from "firebase/functions";
import firebase from "./firebase";

const getClaims = async (user) => {
    const idTokenResult = await user.getIdTokenResult();
    return idTokenResult.claims["https://hasura.io/jwt/claims"];
}

const UserContext = react.createContext({user: null, loading: true, token: ""});

const useRecruiter = () => {
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
                if (!claims) {
                    setState(s => {
                        return {...s, user, loading: true}
                    });

                    const functions = getFunctions(firebase, "europe-west2");
                    const changeMeta = httpsCallable(functions, "changeMeta");

                    await changeMeta({user_type: "recruiter"});
                    const token = await user.getIdToken(true);
                    const claims = await getClaims(user);

                    setState(s => {
                        return {...s, user, claims, token, loading: false}
                    });
                    return
                }

                const token = await user.getIdToken();
                setState(s => {
                    return {...s, loading: false, user, claims, token}
                });
                return;
            }

            console.log('user', user)
            setState(s => {
                return {...s, loading: false, user}
            });
        });

        return () => {
            listen();
        }
    }, [auth]);

    return state;
}

const RecruiterAuthProvider = ({children}) => {
    const state = useRecruiter();

    return (
        <UserContext.Provider value={state}>
            {children}
        </UserContext.Provider>
    );
}

const useCandidate = () => {
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
                setState(s => {
                    return {...s, loading: false, user, claims, token}
                });
                return;
            }

            setState(s => {
                return {...s, loading: false, user, token: "", claims: null}
            });
        });

        return () => {
            listen();
        }
    }, [auth]);

    return state;
}

const CandidateAuthProvider = ({children}) => {
    const state = useCandidate();

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


export {CandidateAuthProvider, RecruiterAuthProvider, useFirebaseAuth, getClaims};