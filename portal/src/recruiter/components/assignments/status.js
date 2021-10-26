import React from "react";

const AssignmentStatus = (props) => {
    let colour = "orange";

    switch (props.status) {
        case "unsent":
            colour = "gray";
            break;
        case "sending":
            colour = "gray";
            break;
        case "sent":
            colour = "blue";
            break;
        case "scheduled":
            colour = "pink"
            break;
        case "in progress":
            colour = "yellow";
            break;
        case "submitted":
            colour = "green";
            break;
        case "missed":
            colour = "red";
            break;
        default:
            colour = "orange"
    }

    return (
        <span
            className={"relative inline-block px-3 py-1 font-semibold text-" + colour + "-900 leading-tight"}>
            <span aria-hidden className={"absolute inset-0 bg-" + colour + "-200 opacity-50 rounded-full"}></span>
            <span className="relative">{props.status}</span>
        </span>
    )
}

export default AssignmentStatus;
