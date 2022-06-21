import React from "react";

type CardProps = {
    title?: string;
    children?: React.ReactNode;
    width?: number;
};

export class Card extends React.Component<CardProps> {
    render() {
        let title;
        if (this.props.title) {
            title = <div className="card-title">{this.props.title}</div>;
        }
        let style: any = {};
        if (this.props.width) {
            style.width = `${this.props.width}rem`;
            style.maxWidth = `${this.props.width}rem`;
        }
        return (
            <div className="card shadow-xl bg-base-100 h-min" style={style}>
                <div className="card-body flex flex-col p-4 items-center">
                    {title}
                    <div className="w-full flex flex-col gap-y-4 p-4">
                        {this.props.children}
                    </div>
                </div>
            </div>
        );
    }
}
