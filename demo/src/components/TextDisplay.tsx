import React from "react";

type TextDisplayProps = {
    label?: string;
    inlineLabel?: string;
    text?: string;
};

export class TextDisplay extends React.Component<TextDisplayProps> {
    render() {
        const text = (
            <div className="input input-bordered input-md flex items-center overflow-x">
                {this.props.text}
            </div>
        );
        if (this.props.inlineLabel)
            return (
                <label className="input-group">
                    <span>{this.props.inlineLabel}</span>
                    {text}
                </label>
            );
        else if (this.props.label) {
            return (
                <div className="form-control w-full max-w-xs">
                    <label className="label">
                        <span className="label-text">{this.props.label}</span>
                    </label>
                    {text}
                </div>
            );
        } else {
            return text;
        }
    }
}
