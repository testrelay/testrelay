import React from "react";

const Languages = (props) => {
    const lang = props.languages.map((e, i) => {
        return <div key={i} class="px-3 py-1 bg-gray-700 rounded-badge text-white">{e.language.name}</div>
    })

    return (
        <div className="flex flex-wrap items-start space-x-2">
            {lang}
        </div>
    )
}

export default Languages;