import React from "react";
import ReactDOM from "react-dom";
import PortalApp from "./recruiter/app";
import CandidateApp from "./candidate/app";
import "./index.css";

let App = PortalApp
const domain =  window.location.host.split('.')[1] ? window.location.host.split('.')[0] : 'app';

if (domain === "candidates") {
    App = CandidateApp
}

ReactDOM.render(
    <App/>
    ,
    document.getElementById("root")
);
