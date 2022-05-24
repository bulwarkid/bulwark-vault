import React from "react";

type KeyDirectoryDisplayProps = {
  keyDirectory?: string;
};

export class KeyDirectoryDisplay extends React.Component<KeyDirectoryDisplayProps> {
  render() {
    return (
      <div className="card shadow-xl bg-base-100 w-[32rem] mw-[32rem]">
        <div className="card-body items-center">
          <div className="card-title text-center">Key Directory</div>
          <div className="max-w-full overflow-x-auto bg-neutral text-neutral-content rounded-lg mx-4 my-4 px-4 py-2">
            <pre>{this.props.keyDirectory}</pre>
          </div>
        </div>
      </div>
    );
  }
}
