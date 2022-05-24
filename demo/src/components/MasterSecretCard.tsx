import React from "react";

const letters = "abcdefghijklmnopqrstuvwxyz";
const lowercase = letters.split("");
const uppercase = letters.toUpperCase().split("");
const numbers = "0123456789".split("");

type ColorCodedBase64Props = {
  text: string;
};

class ColorCodedBase64 extends React.Component<ColorCodedBase64Props> {
  render() {
    const textDivs = [];
    let i = 0;
    for (const letter of this.props.text) {
      let colorClass;
      if (lowercase.includes(letter)) {
        colorClass = "text-purple-400";
      } else if (uppercase.includes(letter)) {
        colorClass = "text-sky-400";
      } else if (numbers.includes(letter)) {
        colorClass = "text-green-400";
      } else {
        colorClass = "text-slate-200";
      }
      textDivs.push(
        <div key={i} className={"inline " + colorClass}>
          {letter}
        </div>
      );
      i++;
    }
    return <span>{textDivs}</span>;
  }
}

type MasterSecretCardProps = {
  masterSecret: string;
  onLogout?: () => void;
};

export class MasterSecretCard extends React.Component<MasterSecretCardProps> {
  render() {
    return (
      <div className="card shadow-xl bg-base-100">
        <div className="card-body items-center">
          <div className="card-title text-l">Master Secret</div>
          <div className="text-l bg-neutral rounded px-2 py-1">
            <ColorCodedBase64 text={this.props.masterSecret} />
          </div>
          <div className="card-actions justify-end">
            <div
              className="btn btn-error btn-sm"
              onClick={() => {
                if (this.props.onLogout) {
                  this.props.onLogout();
                }
              }}
            >
              Log Out
            </div>
          </div>
        </div>
      </div>
    );
  }
}
