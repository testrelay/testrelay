import React from "react";

const Languages = (props) => {
    const lang = props.languages.map((e, i) => {
        return <div key={i} class="badge badge-lg">{e.language.name}</div>
    })

    return (
        <div className="flex flex-wrap items-start space-x-2">
            {lang}
        </div>
    )
}

export default Languages;