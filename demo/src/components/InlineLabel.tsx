import React from "react";

type Props = {
    label: string;
    children?: React.ReactNode;
};

export class InlineLabel extends React.Component<Props> {
    render() {
        return (
            <div className="form-control w-full max-w-lg">
                <label className="input-group whitespace-nowrap">
                    <span>{this.props.label}</span>
                    {this.props.children}
                </label>
            </div>
        );
    }
}
