import React from "react";

type TextInputProps = {
    label?: string;
    placeholder: string;
    type?: string;
};

export class TextInput extends React.Component<TextInputProps> {
    data: string;
    constructor(props: TextInputProps) {
        super(props);
        this.data = "";
    }
    render() {
        const input = (
            <input
                type={this.props.type ?? "text"}
                placeholder={this.props.placeholder}
                className="input input-bordered w-full max-w-s"
                onChange={(e) => {
                    this.data = e.target.value;
                }}
            />
        );
        if (this.props.label) {
            return (
                <div className="form-control w-full max-w-xs">
                    <label className="label">
                        <span className="label-text">{this.props.label}</span>
                    </label>
                    {input}
                </div>
            );
        } else {
            return input;
        }
    }
}
