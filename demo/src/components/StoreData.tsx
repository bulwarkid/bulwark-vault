import React from "react";
import { Label } from "./Label";
import { TextInput } from "./TextInput";

type StoreDataProps = {
    onSubmit?: (path: string, data: string) => void;
};

export class StoreData extends React.Component<StoreDataProps> {
    pathRef = React.createRef<TextInput>();
    dataRef = React.createRef<TextInput>();

    render() {
        return (
            <div className="card shadow-xl bg-base-100 w-max h-max">
                <div className="card-body items-center">
                    <div className="card-title">Store Data</div>
                    <Label label="Path">
                        <TextInput
                            ref={this.pathRef}
                            placeholder="/path/to/data"
                        />
                    </Label>
                    <Label label="Data">
                        <TextInput ref={this.dataRef} placeholder="Data" />
                    </Label>
                    <button
                        className="btn btn-primary"
                        onClick={() => {
                            if (this.props.onSubmit) {
                                this.props.onSubmit(
                                    this.pathRef.current?.data || "",
                                    this.dataRef.current?.data || ""
                                );
                            }
                        }}
                    >
                        Store
                    </button>
                </div>
            </div>
        );
    }
}
