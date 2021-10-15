import { React } from "react";
import { Link } from "react-router-dom";

const ListTitle = (props) => {
    return (
        <div className="py-4 border-b-2 mb-8 border-primary flex">
            <div className="flex-1 flex flex-col justify-center">
                <h2 className="text-xl text-primary">{props.title}</h2>
            </div>
            <div>
                <Link to={props.link} className="btn btn-primary w-auto">{props.button}</Link>
            </div>
        </div>
    )
}

export default ListTitle;