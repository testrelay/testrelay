import React from "react";
import { ListTitle } from "../../../components";
import List from "../../components/users/list";

const UserList = (props) => {
    return (
        <div>
            <ListTitle link="/users/create" button="Invite User" title="Users" />
            <List />
        </div>
    )
}

export default UserList;