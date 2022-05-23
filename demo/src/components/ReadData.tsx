import React from "react";
import { TextInput } from "./TextInput";

type ReadDataState = {
  data: string | undefined;
};

export class ReadData extends React.Component<{}, ReadDataState> {
  constructor(props: {}) {
    super(props);
    this.state = { data: undefined };
  }
  render() {
    return (
      <div className="card shadow-xl bg-base-100">
        <div className="card-body">
          <div className="form-control">
            <div className="input-group">
              <TextInput placeholder="/path/to/data" />
              <button className="btn">Read</button>
            </div>
          </div>
          <div className="bg-neutral text-neutral-content rounded-lg px-4 py-2 w-[512px] h-[256px]">
            {this.state.data}
          </div>
        </div>
      </div>
    );
  }
}
