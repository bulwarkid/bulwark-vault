import React from "react";

type KeyDirectoryDisplayProps = {
  keyDirectory?: string;
};

export class KeyDirectoryDisplay extends React.Component<KeyDirectoryDisplayProps> {
  render() {
    return (
      <div className="card shadow-xl bg-base-100 w-[32rem]">
        <div className="card-body items-center">
          <div className="card-title text-center">Key Directory</div>
          {this.props.keyDirectory}
        </div>
      </div>
    );
  }
}
