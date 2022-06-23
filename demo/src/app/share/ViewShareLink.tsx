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

        const { publicKey, encryptionKey } = this.getKeys(
            this.props.rawFragment
        );
        if (!publicKey || !encryptionKey) {
            return emptyBox;
        }

        return (
            <div className="flex flex-col gap-y-4 w-full">
                <InlineLabel label="Public Key">
                    <TextDisplay text={publicKey} />
                </InlineLabel>
                <InlineLabel label="Encryption Key">
                    <TextDisplay text={encryptionKey} />
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
        const { publicKey, encryptionKey } = this.getKeys(props.rawFragment);
        setImmediate(() => {
            getAuthData(publicKey, encryptionKey).then((data) => {
                this.setState({ data });
            });
        });
    };

    getKeys = (rawFragment: string) => {
        const parts = rawFragment.split("~");
        if (parts.length === 2) {
            const publicKey = parts[0];
            const encryptionKey = parts[1];
            return { publicKey, encryptionKey };
        } else if (parts.length === 3) {
            const publicKey = parts[0];
            const encryptionKey = parts[2];
            return { publicKey, encryptionKey };
        } else {
            return { publicKey: undefined, privateKey: undefined };
        }
    };
}
