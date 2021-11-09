import { initializeApp } from "firebase/app";
// import { getAnalytics } from "firebase/analytics";

const firebaseConfig = {
    apiKey: process.env.REACT_APP_FIREBASE_API_KEY,
    authDomain: process.env.REACT_APP_FIREBASE_API_AUTH_DOMAIN,
    projectId: process.env.REACT_APP_FIREBASE_API_PROJECT_ID,
    storageBucket: process.env.REACT_APP_FIREBASE_API_STORAGE_BUCKET,
    messagingSenderId: process.env.REACT_APP_FIREBASE_API_MESSAGING_SENDER,
    appId: process.env.REACT_APP_FIREBASE_API_APP_ID,
    measurementId: process.env.REACT_APP_FIREBASE_API_MEASUREMENT_ID,
    databaseURL: process.env.REACT_APP_FIREBASE_DATABASE,
};

// Initialize Firebase
export default initializeApp(firebaseConfig);