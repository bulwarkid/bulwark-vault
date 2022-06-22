import React from "react";
import { createAuthData } from "../wasm/vault";
import { TextArea } from "../components/TextArea";
import { TextDisplay } from "../components/TextDisplay";
import { TextInput } from "../components/TextInput";

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
                    <TextDisplay
                        inlineLabel="Public Key"
                        text={this.state.publicKey}
                    />
                    <TextDisplay
                        inlineLabel="Private Key"
                        text={this.state.privateKey}
                    />
                    <TextDisplay
                        inlineLabel="Encryption Key"
                        text={this.state.encryptionKey}
                    />
                    <div className="form-control max-w-full">
                        <div className="input-group max-w-full whitespace-nowrap">
                            <span>Share Link</span>
                            <TextDisplay text={this.state.shareLink} />
                            <button
                                className="btn"
                                onClick={() =>
                                    this.copyLink(this.state.shareLink)
                                }
                            >
                                Copy
                            </button>
                        </div>
                    </div>
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
