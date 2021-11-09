import React from "react";


const Pagination = (props) => {
    if (props.limit > props.total) {
        return (null)
    }

    const pages = Math.ceil(props.total / props.limit);

    const btns = [];
    for (let i = 0; i < pages; i++) {
        if (i === props.page) {
            btns.push((<button key={i}
                               className="z-10 bg-indigo-50 border-indigo-500 text-indigo-600 relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                               value={i}>{i + 1}</button>))
            continue
        }

        btns.push((
            <button key={i}
                    className="bg-white border-gray-300 text-gray-500 hover:bg-gray-50 relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                    value={i} onClick={() => {
                props.setState(i)
            }}>{i + 1}</button>))
    }

    let previous = (
        <button
            className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
            onClick={() => {
                if (props.page !== 0) {
                    props.setState(props.page - 1)
                }
            }}>
            <span className="sr-only">Previous</span>
            <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                 aria-hidden="true">
                <path fillRule="evenodd"
                      d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
                      clipRule="evenodd"/>
            </svg>
        </button>
    )

    let next = (
        <button
            className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
            onClick={() => {
                if (props.page !== pages - 1) {
                    props.setState(props.page + 1)
                }
            }}>
            <span className="sr-only">Next</span>
            <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                 aria-hidden="true">
                <path fillRule="evenodd"
                      d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                      clipRule="evenodd"/>
            </svg>
        </button>
    )

    return (
        <div className="flex justify-center py-4">
            <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                {previous}
                {btns}
                {next}
            </nav>
        </div>
    )
}

export default Pagination;