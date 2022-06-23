import React from "react";
import { createAuthData } from "../../wasm/vault";
import { TextArea } from "../../components/TextArea";
import { TextDisplay } from "../../components/TextDisplay";
import { TextInput } from "../../components/TextInput";
import { InlineLabel } from "../../components/InlineLabel";

type CreateShareLinkState = {
    loading?: boolean;
    publicKey?: string;
    privateKey?: string;
    encryptionKey?: string;
    shareLink?: string;
    editLink?: string;
};

export class CreateShareLink extends React.Component<{}, CreateShareLinkState> {
    dataRef = React.createRef<TextInput>();
    constructor(props: {}) {
        super(props);
        this.state = {};
    }
    render() {
        let output;
        if (this.state.loading) {
            output = <div className="btn btn-ghost loading" />;
        } else if (this.state.publicKey) {
            output = (
                <>
                    <InlineLabel label="Public Key">
                        <TextDisplay text={this.state.publicKey} />
                    </InlineLabel>
                    <InlineLabel label="Private Key">
                        <TextDisplay text={this.state.privateKey} />
                    </InlineLabel>
                    <InlineLabel label="Encryption Key">
                        <TextDisplay text={this.state.encryptionKey} />
                    </InlineLabel>
                    <InlineLabel label="Share Link">
                        <TextDisplay text={this.state.shareLink} />
                        <button
                            className="btn"
                            onClick={() => this.copyLink(this.state.shareLink)}
                        >
                            Copy
                        </button>
                    </InlineLabel>
                    <InlineLabel label="Edit Link">
                        <TextDisplay text={this.state.editLink} />
                        <button
                            className="btn"
                            onClick={() => this.copyLink(this.state.editLink)}
                        >
                            Copy
                        </button>
                    </InlineLabel>
                </>
            );
        }
        return (
            <div>
                <div className="flex flex-col items-center gap-y-4">
                    <TextArea ref={this.dataRef} placeholder="Data to Share" />
                    <button
                        className="btn max-w-sm"
                        onClick={this.createShareLink}
                    >
                        Create Link
                    </button>
                </div>
                {output}
            </div>
        );
    }

    createShareLink = async () => {
        this.setState({ loading: true });
        const data = this.dataRef.current?.data;
        if (!data) {
            return;
        }
        const keys = await createAuthData(data);
        const shareLink = `http://localhost:3000/share#${keys.publicKey}~${keys.encryptionKey}`;
        const editLink = `http://localhost:3000/share#${keys.publicKey}~${keys.privateKey}~${keys.encryptionKey}`;
        this.setState({
            publicKey: keys.publicKey,
            privateKey: keys.privateKey,
            encryptionKey: keys.encryptionKey,
            shareLink,
            editLink,
            loading: false,
        });
    };

    copyLink = (link?: string) => {
        if (!link) {
            return;
        }
        navigator.clipboard.writeText(link);
    };
}
