import React from "react";

type Props = {
    label: string;
    children?: React.ReactNode;
};

export class Label extends React.Component<Props> {
    render() {
        return (
            <div className="form-control w-full max-w-xs">
                <label className="label">
                    <span className="label-text">{this.props.label}</span>
                </label>
                {this.props.children}
            </div>
        );
    }
}
