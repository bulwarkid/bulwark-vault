import React from "react";
import ReactDOM from "react-dom/client";

type AlertProps = {
  text: string;
};

class Alert extends React.Component<AlertProps> {
  render() {
    return (
      <div className="absolute top-48 left-1/2 alert shadow-lg alert-error max-w-xl -translate-x-1/2 -translate-y-1/2 ">
        {this.props.text}
      </div>
    );
  }
}

export function showAlert(text: string) {
  const alertNode = window.document.createElement("div");
  ReactDOM.createRoot(alertNode).render(<Alert text={text} />);
  window.document.body.append(alertNode);
  setTimeout(() => {
    window.document.body.removeChild(alertNode);
  }, 2000);
}
