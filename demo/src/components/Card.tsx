import React from "react";

type CardProps = {
    title?: string;
    children?: React.ReactNode;
};

export class Card extends React.Component<CardProps> {
    render() {
        let title;
        if (this.props.title) {
            title = <div className="card-title">{this.props.title}</div>;
        }
        return (
            <div className="card shadow-xl bg-base-100 w-[32rem] mw-[32rem] h-min">
                <div className="card-body  flex flex-col p-4 items-center">
                    {title}
                    <div className="w-full flex flex-col gap-y-4 p-4">
                        {this.props.children}
                    </div>
                </div>
            </div>
        );
    }
}
