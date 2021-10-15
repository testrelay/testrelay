const formatDate = (date) => {
    let d = new Date(date),
        month = '' + (d.getMonth() + 1),
        day = '' + d.getDate(),
        year = d.getFullYear();

    if (month.length < 2)
        month = '0' + month;
    if (day.length < 2)
        day = '0' + day;

    return [year, month, day].join('-');
}

const dateToReadable = (date) => {
    const monthNames = ["January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ];

    return monthNames[date.getMonth()] + " " + date.getDate() + ", " + date.getFullYear();
}

const testLimitToReadable = (timeLimit) => {
    const hours = Math.floor(timeLimit / 3600);

    if (hours > 25) {
        const days = Math.floor(hours / 24)

        if (days > 7) {
            const weeks = Math.floor(days / 7)

            return weeks > 1 ? weeks + " weeks" : weeks + " week";
        }

        return days + " days";
    }

    return hours > 1 ? hours + " hours" : hours + " hour";
}


export { formatDate, dateToReadable, testLimitToReadable };