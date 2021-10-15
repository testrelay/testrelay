import { React } from "react";
import { Link } from "react-router-dom";

const ListTitle = (props) => {
    return (
        <div className="py-2 border-b-4 mb-6 flex">
            <div className="flex-1 flex flex-col justify-center">
                <h2 className="text-xl font-bold">{props.title}</h2>
            </div>
            <div>
                <Link to={props.link} className="btn btn-primary bg-indigo-600 w-auto h-10 min-h-0">{props.button}</Link>
            </div>
        </div>
    )
}

export default ListTitle;