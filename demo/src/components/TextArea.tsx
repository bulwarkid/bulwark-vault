import React from "react";

type TextAreaProps = {
    label?: string;
    placeholder: string;
    type?: string;
    initialData?: string;
};

type TextAreaState = {
    data?: string;
};

export class TextArea extends React.Component<TextAreaProps, TextAreaState> {
    constructor(props: TextAreaProps) {
        super(props);
        this.state = { data: this.props.initialData };
    }
    render() {
        return (
            <textarea
                placeholder={this.props.placeholder}
                className="textarea textarea-bordered w-full max-w-lg"
                value={this.state.data}
                onChange={(e) => {
                    this.setState({ data: e.target.value });
                }}
            />
        );
    }
}
