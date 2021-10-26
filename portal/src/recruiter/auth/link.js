import {
    fetchSignInMethodsForEmail,
    signInWithEmailAndPassword,
    signInWithPopup,
    EmailAuthProvider,
    GoogleAuthProvider,
    linkWithCredential,
} from "firebase/auth";

const handleAuthError = async (auth, provider, error) => {
    if (error.code === 'auth/account-exists-with-different-credential') {
        const existingEmail = error.customData.email;
        const previousCredential = provider.credentialFromError(error);

        const providers = await fetchSignInMethodsForEmail(auth, existingEmail);
        if (providers.indexOf(EmailAuthProvider.PROVIDER_ID) !== -1) {
            const password = window.prompt('Please provide the password for ' + existingEmail);

            try {
                await signInWithEmailAndPassword(auth, existingEmail, password);
            } catch (error) {
                return { error: errorToReadable(error) }
            }
        } else if (providers.indexOf(GoogleAuthProvider.PROVIDER_ID) !== -1) {
            const google = new GoogleAuthProvider();
            google.setCustomParameters({ 'login_hint': existingEmail });

            try {
                await signInWithPopup(auth, google);
            } catch (error) {
                return { creds: null, error: errorToReadable(error) }
            }
        }

        if (auth.currentUser) {
            try {
                const creds = await linkWithCredential(auth.currentUser, previousCredential)

                return { creds }
            } catch (error) {
                return { error: errorToReadable(error) }
            }
        }

    }


    return { error: errorToReadable(error) }
}

const errorToReadable = (error) => {
    console.log(error);
    return error.message;
}

export { handleAuthError }