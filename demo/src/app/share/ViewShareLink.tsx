import React from "react";
import { getAuthData } from "../../wasm/vault";
import { setImmediate } from "../../util";
import { TextDisplay } from "../../components/TextDisplay";
import { Card } from "../../components/Card";
import { Label } from "../../components/Label";
import { InlineLabel } from "../../components/InlineLabel";

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
            <div className="flex flex-col gap-y-4 w-full">
                <InlineLabel label="Public Key">
                    <TextDisplay text={publicKeyBase64} />
                </InlineLabel>
                <InlineLabel label="Encryption Key">
                    <TextDisplay text={encryptionKeyBase64} />
                </InlineLabel>
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
