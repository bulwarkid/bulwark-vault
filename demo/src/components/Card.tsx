import React from "react";

type CardProps = {
    children?: React.ReactNode;
};

export class Card extends React.Component<CardProps> {
    render() {
        return (
            <div className="card shadow-xl bg-base-100 w-[32rem] mw-[32rem] flex flex-col gap-y-4 p-4">
                {this.props.children}
            </div>
        );
    }
}
