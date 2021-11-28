import { React } from "react";
import { Link } from "react-router-dom";

const ListTitle = (props) => {
    return (
        <div className="py-2 border-b-4 mb-6 flex">
            <div className="flex-1 flex flex-col justify-center">
                <h2 className="text-xl font-bold">{props.title}</h2>
            </div>
            <div>
                <Link to={props.link} className="hidden sm:block hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-3 py-2 w-auto">{props.button}</Link>
                <Link to={props.link} className="sm:hidden hover:bg-indigo-500 bg-indigo-600 text-white text-sm rounded px-2 py-1 w-10 flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 inline" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clipRule="evenodd" />
                    </svg>
                </Link>
            </div>
        </div>
    )
}

export default ListTitle;