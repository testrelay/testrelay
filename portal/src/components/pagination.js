import React from "react";


const Pagination = (props) => {
    if (props.limit > props.total) {
        return (null)
    }

    const pages = Math.ceil(props.total / props.limit);

    const btns = [];
    for (let i = 0; i < pages; i++) {
        if (i === props.page) {
            btns.push((<button key={i} class="btn btn-active" value={i}>{i + 1}</button>))
            continue
        }

        btns.push((<button key={i} class="btn bg-white border-0 text-gray-800 hover:text-white" value={i} onClick={() => { props.setState(i) }}>{i + 1}</button>))
    }

    let previous = (<button class="btn bg-white text-gray-800 border-0 hover:text-white" onClick={() => { props.setState(props.page - 1) }}>Previous</button>)

    if (props.page === 0) {
        previous = <button class="btn btn-disabled">Previous</button>
    }

    let next = (<button class="btn bg-white text-gray-800 hover:text-white border-0" onClick={() => { props.setState(props.page + 1) }}>Next</button>)
    if (props.page === pages - 1) {
        next = <button class="btn btn-disabled">Next</button>
    }

    return (
        <div className="flex items-center justify-center mt-8">
            <div class="btn-group shadow-md">
                {previous}
                {btns}
                {next}
            </div>

        </div>
    )
}

export default Pagination;