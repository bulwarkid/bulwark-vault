import React from "react";
import { CreateShareLink } from "../components/CreateShareLink";
import { ViewShareLink } from "../components/ViewShareLink";

export class ShareData extends React.Component {
    render() {
        let fragment = window.location.hash;
        console.log("Fragment:", fragment);
        if (fragment) {
            fragment = fragment.substring(1); // Remove '#' from start
        }
        return (
            <div className="flex gap-8">
                <CreateShareLink />
                <ViewShareLink rawFragment={fragment} />
            </div>
        );
    }
}
