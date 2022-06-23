import React from "react";

type TextInputProps = {
    placeholder: string;
    type?: string;
    onChange?: () => void;
};

export class TextInput extends React.Component<TextInputProps> {
    data: string;
    constructor(props: TextInputProps) {
        super(props);
        this.data = "";
    }
    render() {
        return (
            <input
                type={this.props.type ?? "text"}
                placeholder={this.props.placeholder}
                className="input input-bordered w-full max-w-sm"
                onChange={(e) => {
                    this.data = e.target.value;
                    if (this.props.onChange) {
                        this.props.onChange();
                    }
                }}
            />
        );
    }
}
