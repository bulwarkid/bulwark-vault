import React from "react";
import { TextInput } from "../../components/TextInput";

type ReadDataProps = {
    onRead?: (path: string) => Promise<string | undefined>;
};

type ReadDataState = {
    data: string | undefined;
};

export class ReadData extends React.Component<ReadDataProps, ReadDataState> {
    pathInputRef = React.createRef<TextInput>();
    constructor(props: {}) {
        super(props);
        this.state = { data: undefined };
    }
    render() {
        return (
            <div className="card shadow-xl bg-base-100">
                <div className="card-body flex flex-col items-center">
                    <div className="card-title">Read Data</div>
                    <div className="form-control w-full">
                        <div className="input-group">
                            <TextInput
                                ref={this.pathInputRef}
                                placeholder="/path/to/data"
                            />
                            <button className="btn" onClick={this.onRead}>
                                Read
                            </button>
                        </div>
                    </div>
                    <div className="bg-neutral text-neutral-content rounded-lg px-4 py-2 w-[512px] h-[256px]">
                        {this.state.data}
                    </div>
                </div>
            </div>
        );
    }

    onRead = async () => {
        const path = this.pathInputRef.current?.data;
        if (path && this.props.onRead) {
            const data = await this.props.onRead(path);
            this.setState({ data });
        }
    };
}
