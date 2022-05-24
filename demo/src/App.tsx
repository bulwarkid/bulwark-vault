import React from "react";
import { LoginCard } from "./components/LoginCard";
import { MasterSecretCard } from "./components/MasterSecretCard";
import { ReadData } from "./components/ReadData";
import { StoreData } from "./components/StoreData";
import * as vault from "./wasm/vault";

type AppState = {
  username?: string;
  password?: string;
};

export class App extends React.Component<{}, AppState> {
  constructor(props: {}) {
    super(props);
    this.state = {};
  }
  render() {
    let content;
    if (!this.state.username || !this.state.password) {
      content = <LoginCard onLogin={this.onLogin} />;
    } else {
      content = (
        <div className="flex space-y-4 flex-col items-center">
          <MasterSecretCard
            masterSecret={vault.getMasterSecret()}
            onLogout={this.onLogout}
          />
          <StoreData />
        </div>
      );
    }
    return (
      <div className="w-full h-full flex justify-center items-center flex-col">
        {content}
      </div>
    );
  }

  onLogin = (username: string | null, password: string | null) => {
    if (!username || !password) {
      return;
    }
    vault.loginToVault(username, password).then(() => {
      this.setState({ username, password });
    });
  };

  onLogout = () => {
    this.setState({ username: undefined, password: undefined });
  };
}
