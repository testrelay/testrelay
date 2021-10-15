import { React } from "react";


import CreateTestForm from "../../components/tests/create-test-form";

const CreateTest = (props) => {
	return (
		<div>
			<div className="py-4 border-b-4 mb-6">
				<h2 className="text-xl font-bold">Create Test</h2>
			</div>
			<CreateTestForm business={props.business} />
		</div>
	)
}

export default CreateTest;
