import React from "react";
import { LoginCard } from "./LoginCard";
import { MasterSecretCard } from "./MasterSecretCard";
import { ReadData } from "./ReadData";
import { StoreData } from "./StoreData";

import "./style/App.css";

export class App extends React.Component {
  render() {
    return (
      <div className="w-full h-full flex justify-center items-center flex-col">
        <LoginCard />
      </div>
    );
  }
}
