import React from "react";
import { PersonalVault } from "./PersonalVault";
import { ShareData } from "./ShareData";

export class App extends React.Component {
    render() {
        let content;
        if (window.location.pathname.startsWith("/share")) {
            content = <ShareData />;
        } else {
            content = <PersonalVault />;
        }
        return (
            <div className="w-full h-full flex justify-center items-center flex-col">
                <div className="w-full h-32 absolute top-8 left-8">
                    <img
                        src="company-name.png"
                        width="194"
                        height="59"
                        alt="BulwarkID Company Brand"
                    ></img>
                </div>
                {content}
            </div>
        );
    }
}
