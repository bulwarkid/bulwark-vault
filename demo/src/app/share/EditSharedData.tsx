import React from "react";
import { InlineLabel } from "../../components/InlineLabel";
import { Label } from "../../components/Label";
import { TextArea } from "../../components/TextArea";
import { TextDisplay } from "../../components/TextDisplay";
import { TextInput } from "../../components/TextInput";
import { getAuthData } from "../../wasm/vault";

type EditSharedDataProps = {
    initialLink?: string;
};

type EditSharedDataState = {
    link?: string;
    data?: string;
};

export class EditSharedData extends React.Component<
    EditSharedDataProps,
    EditSharedDataState
> {
    linkRef = React.createRef<TextInput>();
    constructor(props: EditSharedDataProps) {
        super(props);
        this.state = { link: props.initialLink };
    }

    componentDidUpdate(prevProps: EditSharedDataProps) {
        if (this.props.initialLink !== prevProps.initialLink) {
            this.updateLink(this.props.initialLink);
        }
    }

    render() {
        const editLink = this.state.link;
        let publicKey, privateKey, encryptionKey;
        if (editLink) {
            const keys = this.getKeys(editLink);
            publicKey = keys.publicKey;
            privateKey = keys.privateKey;
            encryptionKey = keys.encryptionKey;
        }
        const showErrorMessage = !!publicKey && !privateKey;
        return (
            <div className="flex flex-col gap-y-4 items-center">
                {showErrorMessage && (
                    <div className="alert alert-error">
                        {
                            "You need the private key in order to edit. Did you use the Share Link?"
                        }
                    </div>
                )}
                <InlineLabel label="Edit Link">
                    <TextInput
                        ref={this.linkRef}
                        placeholder="Enter edit link..."
                        onChange={this.onLinkUpdate}
                    />
                </InlineLabel>
                <InlineLabel label="Public Key">
                    <TextDisplay text={publicKey} />
                </InlineLabel>
                <InlineLabel label="Private Key">
                    <TextDisplay text={privateKey} />
                </InlineLabel>
                <InlineLabel label="Encryption Key">
                    <TextDisplay text={encryptionKey} />
                </InlineLabel>
                <Label label="Data">
                    <TextArea
                        placeholder="Data..."
                        initialData={this.state.data}
                    />
                </Label>
                <div className="flex justify-center">
                    <button className="btn max-w-sm">Update</button>
                </div>
            </div>
        );
    }

    onLinkUpdate = async () => {
        const link = this.linkRef.current?.data;
        await this.updateLink(link);
    };

    getKeys = (link: string) => {
        const fragment = link.split("#").pop();
        let publicKey, privateKey, encryptionKey;
        if (fragment) {
            const fragmentParts = fragment.split("~");
            if (fragmentParts.length === 3) {
                publicKey = fragmentParts[0];
                privateKey = fragmentParts[1];
                encryptionKey = fragmentParts[2];
            } else if (fragmentParts.length === 2) {
                publicKey = fragmentParts[0];
                encryptionKey = fragmentParts[1];
            }
        }
        return { publicKey, privateKey, encryptionKey };
    };

    updateLink = async (link?: string) => {
        this.setState({ link }, () => {
            if (link) {
                const keys = this.getKeys(link);
                if (keys.publicKey && keys.privateKey && keys.encryptionKey) {
                    this.updateData(
                        keys.publicKey,
                        keys.privateKey,
                        keys.encryptionKey
                    );
                }
            }
        });
    };

    updateData = async (
        publicKey: string,
        privateKey: string,
        encryptionKey: string
    ) => {
        const data = await getAuthData(publicKey, encryptionKey);
        this.setState({ data });
    };
}
