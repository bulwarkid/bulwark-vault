import React from "react";
import { InlineLabel } from "../components/InlineLabel";
import { TextInput } from "../components/TextInput";

export class EditSharedData extends React.Component {
    render() {
        return (
            <div>
                <InlineLabel label="Public Key">
                    <TextInput placeholder="Enter Public Key..."></TextInput>
                </InlineLabel>
            </div>
        );
    }
}
