import React from "react";

type TextDisplayProps = {
    text?: string;
};

export class TextDisplay extends React.Component<TextDisplayProps> {
    render() {
        return (
            <div className="input input-bordered input-md flex items-center overflow-x-scroll whitespace-nowrap">
                {this.props.text}
            </div>
        );
    }
}
