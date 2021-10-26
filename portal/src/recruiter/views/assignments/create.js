import React from "react";
import CreateCandidateForm from "../../components/assignments/create-candidate-form";

const CandidateCreate = () => {
	return (
		<div>
			<div className="py-4 border-b-4 mb-6">
				<h2 className="text-xl font-bold">Create Assignment</h2>
			</div>
			<CreateCandidateForm />
		</div>
	);
}

export default CandidateCreate;
