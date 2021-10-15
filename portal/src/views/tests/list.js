import React from "react";

import TestList from "../../components/tests/list";
import { ListTitle } from "../../components";

const TestListView = () => {
	return (
		<div>
			<ListTitle link="/tests/create" button="Create Test" title="Tests" />
			<TestList />
		</div>
	)
}

export default TestListView;
