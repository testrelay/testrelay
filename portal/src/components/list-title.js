import { React } from "react";
import { Link } from "react-router-dom";

const ListTitle = (props) => {
    return (
        <div className="py-2 border-b-4 mb-6 flex">
            <div className="flex-1 flex flex-col justify-center">
                <h2 className="text-xl font-bold">{props.title}</h2>
            </div>
            <div>
                <Link to={props.link} className="hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-4 py-2 w-auto">{props.button}</Link>
            </div>
        </div>
    )
}

export default ListTitle;