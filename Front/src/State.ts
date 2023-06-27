import React, {useState} from "react";

export class State<T> {
    private readonly state: T;
    private readonly setState: React.Dispatch<React.SetStateAction<T>>;
    public constructor(value: T) {
        [this.state, this.setState] = useState<T>(value);
    }

    public Get(): T {
        return this.state;
    }

    public Set(value: T) {
        this.setState(value);
    }
}