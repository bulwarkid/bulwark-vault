import React from "react";

type TextAreaProps = {
    label?: string;
    placeholder: string;
    type?: string;
    initialData?: string;
};

export class TextArea extends React.Component<TextAreaProps> {
    data: string;
    constructor(props: TextAreaProps) {
        super(props);
        this.data = this.props.initialData ?? "";
    }
    render() {
        return (
            <textarea
                placeholder={this.props.placeholder}
                className="textarea textarea-bordered w-full max-w-md"
                value={this.props.initialData}
                onChange={(e) => {
                    this.data = e.target.value;
                }}
            />
        );
    }
}
