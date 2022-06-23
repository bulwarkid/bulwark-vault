import React from "react";

type TextDisplayProps = {
    text?: string;
    area?: boolean;
};

export class TextDisplay extends React.Component<TextDisplayProps> {
    render() {
        let className = "flex items-center grow bg-base-200";
        if (this.props.area) {
            className += " textarea";
        } else {
            className += " input input-md whitespace-nowrap overflow-x-scroll";
        }
        return <div className={className}>{this.props.text}</div>;
    }
}
