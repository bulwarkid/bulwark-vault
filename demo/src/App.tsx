import React from "react";
import { KeyDirectoryDisplay } from "./components/KeyDirectoryDisplay";
import { LoginCard } from "./components/LoginCard";
import { MasterSecretCard } from "./components/MasterSecretCard";
import { ReadData } from "./components/ReadData";
import { StoreData } from "./components/StoreData";
import * as vault from "./wasm/vault";

type AppState = {
  username?: string;
  password?: string;
  masterSecret?: string;
  keyDirectory?: string;
};

export class App extends React.Component<{}, AppState> {
  constructor(props: {}) {
    super(props);
    this.state = {};
  }
  render() {
    let content;
    if (
      !this.state.username ||
      !this.state.password ||
      !this.state.masterSecret
    ) {
      content = <LoginCard onLogin={this.onLogin} />;
    } else {
      content = (
        <div className="flex space-y-4 flex-col items-center">
          <MasterSecretCard
            masterSecret={this.state.masterSecret}
            onLogout={this.onLogout}
          />
          <div className="flex space-x-4">
            <StoreData onSubmit={this.onStore} />
            <ReadData onRead={this.onRead} />
          </div>
          <KeyDirectoryDisplay keyDirectory={this.state.keyDirectory} />
        </div>
      );
    }
    return (
      <div className="w-full h-full flex justify-center items-center flex-col">
        <div className="w-full h-32 absolute top-8 left-8">
          <img src="company-name.png" width="194" height="59"></img>
        </div>
        {content}
      </div>
    );
  }

  onLogin = async (username: string | null, password: string | null) => {
    if (!username || !password) {
      return;
    }
    await vault.loginToVault(username, password);
    const masterSecret = await vault.getMasterSecret();
    const keyDirectory = await vault.getKeyDirectory();
    this.setState({ username, password, masterSecret, keyDirectory });
  };

  onLogout = () => {
    this.setState({ username: undefined, password: undefined });
  };

  onStore = async (path: string, data: string) => {
    await vault.put(path, data);
    const keyDirectory = await vault.getKeyDirectory();
    this.setState({ keyDirectory });
  };

  onRead = async (path: string) => {
    return await vault.get(path);
  };
}
