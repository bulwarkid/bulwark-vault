import React from "react";

type TextDisplayProps = {
    inlineLabel?: string;
    text?: string;
};

export class TextDisplay extends React.Component<TextDisplayProps> {
    render() {
        const text = (
            <div className="input input-bordered input-md flex items-center overflow-x-scroll whitespace-nowrap">
                {this.props.text}
            </div>
        );
        if (this.props.inlineLabel) {
            return (
                <div className="form-control">
                    <label className="input-group whitespace-nowrap">
                        <span>{this.props.inlineLabel}</span>
                        {text}
                    </label>
                </div>
            );
        } else {
            return text;
        }
    }
}
