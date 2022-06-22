import React from "react";
import { getAuthData } from "../wasm/vault";
import { setImmediate } from "../util";
import { TextDisplay } from "./TextDisplay";
import { Card } from "./Card";
import { Label } from "./Label";

type ViewShareLinkProps = {
    rawFragment: string;
};

type ViewShareLinkState = {
    data?: string;
};

export class ViewShareLink extends React.Component<
    ViewShareLinkProps,
    ViewShareLinkState
> {
    constructor(props: ViewShareLinkProps) {
        super(props);
        this.state = {};
        this.fetchData(props);
    }

    componentDidUpdate(prevProps: ViewShareLinkProps) {
        if (this.props.rawFragment !== prevProps.rawFragment) {
            this.fetchData(this.props);
        }
    }

    render() {
        const emptyBox = (
            <div className="card shadow-xl bg-base-100 w-[32rem] mw-[32rem]">
                No Data found;
            </div>
        );
        if (!this.props.rawFragment) {
            return emptyBox;
        }

        const parts = this.props.rawFragment.split("~");
        if (parts.length !== 2) {
            return emptyBox;
        }
        const publicKeyBase64 = parts[0];
        const encryptionKeyBase64 = parts[1];

        return (
            <div>
                <TextDisplay inlineLabel="Public Key" text={publicKeyBase64} />
                <TextDisplay
                    inlineLabel="Encryption Key"
                    text={encryptionKeyBase64}
                />
                <Label label="Data">
                    <TextDisplay text={this.state.data} />
                </Label>
            </div>
        );
    }

    fetchData = (props: ViewShareLinkProps) => {
        if (!props.rawFragment) {
            return;
        }
        const parts = props.rawFragment.split("~");
        if (parts.length !== 2) {
            return;
        }
        const publicKeyBase64 = parts[0];
        const encryptionKeyBase64 = parts[1];
        setImmediate(() => {
            getAuthData(publicKeyBase64, encryptionKeyBase64).then((data) => {
                this.setState({ data });
            });
        });
    };
}
