const readableError = (code) => {
    switch (code) {
        case "auth/user-not-found":
            return 'no user with these details exists';
        case "auth/wrong-password":
            return 'incorrect password';
        case "auth/popup-closed-by-user":
            return 'oauth popup closed too early, please try again'
        default:
            return "unexpected error, please refresh the browser and try again"
    }
}

export {readableError}