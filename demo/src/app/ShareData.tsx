import React from "react";
import { Card } from "../components/Card";
import { CreateShareLink } from "./CreateShareLink";
import { ViewShareLink } from "./ViewShareLink";

enum Tab {
    SHARE,
    VIEW,
    EDIT,
}
const TAB_LABELS = new Map([
    [Tab.SHARE, "Share Data"],
    [Tab.VIEW, "View Data"],
    [Tab.EDIT, "Edit Data"],
]);

type State = {
    tab: Tab;
};

export class ShareData extends React.Component<{}, State> {
    constructor(props: {}) {
        super(props);
        this.state = { tab: Tab.SHARE };
    }
    render() {
        let fragment = window.location.hash;
        if (fragment) {
            fragment = fragment.substring(1); // Remove '#' from start
        }
        const tabs = [];
        for (const tab of [Tab.SHARE, Tab.VIEW, Tab.EDIT]) {
            let classList = "tab tab-lifted font-bold";
            if (tab === this.state.tab) {
                classList += " tab-active";
            }
            tabs.push(
                <div
                    className={classList}
                    onClick={() => this.setState({ tab })}
                >
                    {TAB_LABELS.get(tab)}
                </div>
            );
        }
        let cardContents;
        if (this.state.tab === Tab.SHARE) {
            cardContents = <CreateShareLink />;
        } else if (this.state.tab === Tab.VIEW) {
            cardContents = <ViewShareLink rawFragment={fragment} />;
        } else if (this.state.tab === Tab.EDIT) {
            cardContents = <div />;
        }
        return (
            <div className="flex gap-8">
                <div>
                    <div className="tabs -mb-px z-10 relative">{tabs}</div>
                    <Card width={48} extraClasses="rounded-tl-none">
                        {cardContents}
                    </Card>
                </div>
            </div>
        );
    }
}
