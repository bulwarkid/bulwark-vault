import React from "react";
import { LoginCard } from "./LoginCard";
import { MasterSecretCard } from "./MasterSecretCard";
import { ReadData } from "./ReadData";
import { StoreData } from "./StoreData";

import "./style/App.css";

function App() {
  return (
    <div className="w-full h-full flex justify-center items-center flex-col">
      <ReadData />
    </div>
  );
}

export default App;
