import React, {useRef, useState} from "react";
import {Link} from "react-router-dom";
import {
    createUserWithEmailAndPassword,
    getAuth,
    GithubAuthProvider,
    GoogleAuthProvider,
    signInWithPopup
} from "firebase/auth";
import firebase from "../../auth/firebase";
import {AlertError} from "../../components/alerts";
import {handleAuthError} from "../auth/link";
import {Loading} from "../../components";

const Register = (props) => {
    const github = new GithubAuthProvider();
    const google = new GoogleAuthProvider();

    const auth = getAuth(firebase);
    const email = useRef();
    const password = useRef();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
                <div className="max-w-md w-full space-y-8 p-8 shadow-md rounded-lg bg-white">
                    <Loading/>
                </div>
            </div>
        )
    }

    const withLoading = (func) => {
        return async () => {
            setLoading(true);
            await func();
            setLoading(false);
        }
    }

    const signup = withLoading(async () => {
        try {
            await createUserWithEmailAndPassword(auth, email.current.value, password.current.value);
        } catch (error) {
            setError(error.message);
        }
    })


    const signupWithGithub = withLoading(async () => {
        try {
            await signInWithPopup(auth, github)
        } catch (e) {
            const {error} = await handleAuthError(auth, github, e);

            if (error) {
                setError(error);
            }
        }
    })

    const signupWithGoogle = withLoading(async () => {
        try {
            await signInWithPopup(auth, google)
        } catch (e) {
            const {error} = await handleAuthError(auth, google, e);

            if (error) {
                setError(error);
            }
        }
    })

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8 p-8 shadow-md rounded-lg bg-white">
                <div>
                    <div className="flex items-center justify-center mx-auto bg-gray-800 shadow rounded-full"
                         style={{width: 80, height: 80}}>
                        <h1 className="mx-auto font-semibold text-lg text-white">
                            <svg height="40pt" viewBox="0 -48 480 480" width="40pt" fill="#fff"
                                 xmlns="http://www.w3.org/2000/svg">
                                <path
                                    d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0"/>
                                <path
                                    d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0"/>
                            </svg>
                        </h1>
                    </div>
                    <h2 className="mt-4 text-center text-3xl font-bold text-gray-900">
                        Create an account
                    </h2>
                </div>
                {error &&
                <AlertError message={error}/>}
                <div className="mt-8 space-y-2">
                    <input type="hidden" name="remember" value="true"/>
                    <div className="rounded-md shadow-sm -space-y-px">
                        <div>
                            <label htmlFor="email-address" className="sr-only">Email address</label>
                            <input name="email" type="email" autoComplete="email" required
                                   className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                   placeholder="Email address" ref={email}/>
                        </div>
                        <div>
                            <label htmlFor="password" className="sr-only">Password</label>
                            <input name="password" type="password" autoComplete="current-password" required
                                   className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                   placeholder="Password" ref={password}/>
                        </div>
                    </div>

                    <div>
                        <button onClick={signup}
                                className="mb-2 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                                <svg className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400"
                                     xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                                     aria-hidden="true">
                                    <path fillRule="evenodd"
                                          d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                                          clipRule="evenodd"/>
                                </svg>
                            </span>
                            Register
                        </button>
                        <div className="flex flex-row justify-center items-center space-x-2 mb-2">
                            <div className="flex-grow">
                                <div className="border-b-2 border-gray-200"></div>
                            </div>
                            <div className="flex-shrink text-center">
                                <p className="text-sm ">or</p>
                            </div>
                            <div className="flex-grow items-center">
                                <div className="border-b-2 border-gray-200"></div>
                            </div>
                        </div>
                        <button onClick={signupWithGithub}
                                className="mb-2 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                                <svg className="h-5 w-5 text-gray-100 group-hover:text-gray-200" width="20px"
                                     height="20px" viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
                                    <g>
                                        <path fill="currentColor"
                                              d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z"></path>
                                    </g>
                                </svg>
                            </span>
                            Register with github
                        </button>
                        <button onClick={signupWithGoogle}
                                className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                                <svg viewBox="0 0 24 24" width="20" height="20">
                                    <g transform="matrix(1, 0, 0, 1, 27.009001, -39.238998)">
                                        <path fill="#4285F4"
                                              d="M -3.264 51.509 C -3.264 50.719 -3.334 49.969 -3.454 49.239 L -14.754 49.239 L -14.754 53.749 L -8.284 53.749 C -8.574 55.229 -9.424 56.479 -10.684 57.329 L -10.684 60.329 L -6.824 60.329 C -4.564 58.239 -3.264 55.159 -3.264 51.509 Z"/>
                                        <path fill="#34A853"
                                              d="M -14.754 63.239 C -11.514 63.239 -8.804 62.159 -6.824 60.329 L -10.684 57.329 C -11.764 58.049 -13.134 58.489 -14.754 58.489 C -17.884 58.489 -20.534 56.379 -21.484 53.529 L -25.464 53.529 L -25.464 56.619 C -23.494 60.539 -19.444 63.239 -14.754 63.239 Z"/>
                                        <path fill="#FBBC05"
                                              d="M -21.484 53.529 C -21.734 52.809 -21.864 52.039 -21.864 51.239 C -21.864 50.439 -21.724 49.669 -21.484 48.949 L -21.484 45.859 L -25.464 45.859 C -26.284 47.479 -26.754 49.299 -26.754 51.239 C -26.754 53.179 -26.284 54.999 -25.464 56.619 L -21.484 53.529 Z"/>
                                        <path fill="#EA4335"
                                              d="M -14.754 43.989 C -12.984 43.989 -11.404 44.599 -10.154 45.789 L -6.734 42.369 C -8.804 40.429 -11.514 39.239 -14.754 39.239 C -19.444 39.239 -23.494 41.939 -25.464 45.859 L -21.484 48.949 C -20.534 46.099 -17.884 43.989 -14.754 43.989 Z"/>
                                    </g>
                                </svg>
                            </span>
                            Register with google
                        </button>
                        <p className="mt-4 text-center text-sm text-gray-600">
                            <Link to="/login" className="text-indigo-500 hover:text-indigo-500">
                                already have an account?
                            </Link>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Register;