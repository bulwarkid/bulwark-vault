import React from "react";

type TextAreaProps = {
    label?: string;
    placeholder: string;
    type?: string;
};

export class TextArea extends React.Component<TextAreaProps> {
    data: string;
    constructor(props: TextAreaProps) {
        super(props);
        this.data = "";
    }
    render() {
        return (
            <textarea
                placeholder={this.props.placeholder}
                className="textarea textarea-bordered w-full max-w-s"
                onChange={(e) => {
                    this.data = e.target.value;
                }}
            />
        );
    }
}
