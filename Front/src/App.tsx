import React, {useEffect} from "react";
import { State } from "./State";
import Task from "./Task";
import { AddProcess, CPUUpdate, Update } from "./Updates";

export let updates: State<Update[]>;

function parser(result: Update[]) {
    updates.Set(result);
}

export default function App() {
    updates = new State<Update[]>([]);
    const err = new State(false);
    const queues: AddProcess[][] = [[], [], [], []];
    useEffect(()=>{
        fetch("http://localhost:3131/updates", {
            method: "GET",
            redirect: "follow" as RequestRedirect,
            mode: "cors" as RequestMode,
        })
            .then((response) => response.json())
            .then(parser)
            .catch(() => err.Set(true));
    });

    return (<>
        {updates.Get().map((update, i) => {
            if (update.Type == "AddProcess") {
                const proc = update as AddProcess;
                queues[proc.QI].push(proc);
                return <></>;
            }
            else {
                const CPU = update as CPUUpdate;
                queues[CPU.QI].splice(0, 1);
                return <Task key={i} cpu={CPU} queues={queues} />;
            }
        })}
    </>);
}
