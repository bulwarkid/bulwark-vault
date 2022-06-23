import React from "react";
import { TextInput } from "../../components/TextInput";

type LoginCardProps = {
    onLogin?: (username: string | null, password: string | null) => void;
};

export class LoginCard extends React.Component<LoginCardProps> {
    usernameRef = React.createRef<TextInput>();
    passwordRef = React.createRef<TextInput>();

    render() {
        return (
            <div className="card shadow-xl bg-base-100">
                <div className="card-body items-center">
                    <div className="card-title text-center items-center">
                        Log In
                    </div>
                    <TextInput ref={this.usernameRef} placeholder="Username" />
                    <TextInput
                        ref={this.passwordRef}
                        placeholder="Password"
                        type="password"
                    />

                    <button
                        className="btn btn-primary"
                        onClick={() => {
                            if (this.props.onLogin) {
                                this.props.onLogin(
                                    this.usernameRef.current?.data || "",
                                    this.passwordRef.current?.data || ""
                                );
                            }
                        }}
                    >
                        Log In
                    </button>
                </div>
            </div>
        );
    }
}
