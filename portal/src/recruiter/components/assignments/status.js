import React from "react";

const AssignmentStatus = (props) => {
    let colour = "bg-orange-200";

    switch (props.status) {
        case "unsent":
            colour = "bg-gray-200";
            break;
        case "sending":
            colour = "bg-gray-200";
            break;
        case "sent":
            colour = "bg-blue-200";
            break;
        case "scheduled":
            colour = "bg-pink-200"
            break;
        case "in progress":
            colour = "bg-yellow-200";
            break;
        case "submitted":
            colour = "bg-green-200";
            break;
        case "missed":
            colour = "bg-red-200";
            break;
        default:
            colour = "bg-orange-200"
    }

    return (
        <span
            className={"relative inline-block px-4 py-1 font-semibold text-gray-700 text-xs leading-tight"}>
            <span aria-hidden className={"absolute inset-0 " + colour + " opacity-50 rounded-full"}/>
            <span className="relative">{props.status}</span>
        </span>
    )
}

export default AssignmentStatus;
