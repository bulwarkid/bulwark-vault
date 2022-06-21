import React from "react";
import { createAuthData } from "../wasm/vault";
import { Card } from "./Card";
import { TextDisplay } from "./TextDisplay";
import { TextInput } from "./TextInput";

type CreateShareLinkState = {
    publicKey?: string;
    privateKey?: string;
    encryptionKey?: string;
    link?: string;
};

export class CreateShareLink extends React.Component<{}, CreateShareLinkState> {
    dataRef = React.createRef<TextInput>();
    constructor(props: {}) {
        super(props);
        this.state = {};
    }
    render() {
        return (
            <Card title="Create Share Link">
                <div className="form-control">
                    <div className="input-group">
                        <TextInput
                            ref={this.dataRef}
                            type="text"
                            placeholder="Data to Share"
                        />
                        <button className="btn" onClick={this.createShareLink}>
                            Create Link
                        </button>
                    </div>
                </div>
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
                    <div className="input-group max-w-full">
                        <TextDisplay text={this.state.link} />
                        <button
                            className="btn"
                            onClick={() => this.copyLink(this.state.link)}
                        >
                            Copy
                        </button>
                    </div>
                </div>
            </Card>
        );
    }

    createShareLink = async () => {
        const data = this.dataRef.current?.data;
        if (!data) {
            return;
        }
        const keys = await createAuthData(data);
        const link = `http://localhost:3000/share#${keys.publicKey}~${keys.encryptionKey}`;
        this.setState({
            publicKey: keys.publicKey,
            privateKey: keys.privateKey,
            encryptionKey: keys.encryptionKey,
            link,
        });
    };

    copyLink = (link?: string) => {
        if (!link) {
            return;
        }
        navigator.clipboard.writeText(link);
    };
}
