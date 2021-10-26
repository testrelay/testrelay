import { React } from "react";
import { ListTitle } from "../../../components";

import CandidateList from "../../components/assignments/candidate-list";

const CandidateListView = () => {
	return (
		<div>
			<ListTitle link="/assignments/create" button="Create Assignment" title="Assignments" />
			<CandidateList />
		</div>
	)
}

export default CandidateListView;
