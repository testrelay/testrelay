import TimeAgo from "javascript-time-ago";
import React from "react";
import en from 'javascript-time-ago/locale/en'
import { dateToReadable } from "../../../components/date";

TimeAgo.addDefaultLocale(en)

const TimelineIcon = (props) => {
    switch (props.event_type) {
        case "sent":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
            )
        case "viewed":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
            )
        case "scheduled":
        case "rescheduled":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
            )
        case "cancelled":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            )
        case "inprogress":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
            )
        case "submitted":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            )
        case "missed":
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
            )
        default:
            return (
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
            )
    }
}

const TimelineBody = (props) => {
    let title = props.event_type.replace("-", " ").toLowerCase();

    const user_type = props.test.candidate_id === props.user.id ? "candidate" : "user";
    let body = (<span>{user_type + " " + props.user.email + " " + title + " the assignment."}</span>);
    let additional = (null);

    if (props.event_type === "scheduled" || props.event_type === "rescheduled") {
        const scheduled = dateToReadable(new Date(props.meta.test_day_chosen));
        additional = (<span>They're due to take the test at <b>{scheduled + " " + props.meta.test_time_chosen + " (" + props.meta.test_timezone_chosen + ")."}</b></span>)
    }

    const timeAgo = new TimeAgo();
    const d = new Date(props.created_at);
    const date = timeAgo.format(d, "round");

    return (
        <div className="order-1 bg-white rounded-lg shadow-md px-6 py-4 ml-8 flex-grow">
            <h3 className="text-gray-800 text-md capitalize">{title}</h3>
            <h4 className="mb-3 text-gray-300 text-sm">{date + " @ " + d.getHours() + ":" + d.getMinutes()}</h4>
            <p className="text-sm leading-snug tracking-wide text-gray-900 text-opacity-100">{body}{additional}</p>
        </div>
    )
}

const TimelineItem = (props) => {
    return (
        <div className="flex items-center w-full">
            <div className="flex-shrink-0 w-10 h-10 rounded-full bg-indigo-500 inline-flex items-center justify-center text-white relative z-10">
                <TimelineIcon event_type={props.event_type} />
            </div>
            <TimelineBody {...props} />
        </div >
    )
}
const Timeline = (props) => {
    let scheduled = false;
    const items = props.events.map(e => {
        if (e.event_type === "scheduled") {
            if (scheduled) {
                e = Object.assign({}, e, { event_type: "rescheduled" })
            }

            scheduled = true;
        }

        return (
            <TimelineItem {...e} test={props.test} />
        )
    });

    return (
        <div className="mx-auto w-full h-full">
            <div className="mt-10 relative wrap h-full pr-4">
                <div className="h-full w-10 absolute inset-0 flex items-center justify-center">
                    <div className="h-full w-1 bg-gray-200 pointer-events-none"></div>
                </div>
                <div className="space-y-4">
                    {items}
                </div>
            </div>
        </div>
    )
}

export default Timeline;