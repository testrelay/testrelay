import React from "react";
import CreateUser from "../../components/users/create";

const UserCreate = (props) => {
    return (
        <div>
            <div className="py-4 border-b-4 mb-6">
                <h2 className="text-xl font-bold">Create User</h2>
            </div>

            <CreateUser />
        </div>
    )
}

export default UserCreate;