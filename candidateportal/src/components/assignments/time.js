const assignmentLimit = (timeLimit) => {
    const hours = Math.floor(timeLimit / 3600);

    if (hours > 25) {
        const days = Math.floor(hours / 24)

        return days + " days";
    }

    return hours > 1 ? hours + " hours" : hours + " hour";
}

export { assignmentLimit };