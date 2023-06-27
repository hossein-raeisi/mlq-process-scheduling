import React, {useEffect} from "react";
import { State } from "./State";
import { CpuTask } from "./Task";
import { AddProcess, CPUUpdate, Update } from "./Updates";

export let updates: State<Update[]>;

export default function App() {
    updates = new State<Update[]>([]);
    const err = new State(false);
    const queues: AddProcess[][] = [[], [], [], []];
    useEffect(()=>{
        fetch("http://localhost:3131/updates", {
            method: "GET",
            redirect: "follow" as RequestRedirect,
        })
            .then((response) => response.json())
            .then((res) => updates.Set(res))
            .catch((r) => { err.Set(true); console.log(r); });
    }, []);

    return (<>
        <Err err={err} />
        {updates.Get().map((update, i) => {
            if (update.Type == "AddProcess") {
                const proc = update as AddProcess;
                proc.AtStr = Math.floor(proc.AT / 60) + ":" + proc.AT % 60;
                queues[proc.QI].push(proc);
                return <></>;
            }
            else {
                const CPU = update as CPUUpdate;
                CPU.StartStr = Math.floor(CPU.Start / 60) + ":" + CPU.Start % 60;
                CPU.EndStr = Math.floor(CPU.End / 60) + ":" + CPU.End % 60;
                const newQueues = Array<AddProcess[]>(4);
                queues.forEach((queue, i) => {
                    newQueues[i] = Array<AddProcess>();
                    queue.forEach((process) => {
                        newQueues[i].push(Object.assign({}, process));
                    });
                });
                queues[CPU.QI].splice(0, 1);
                return <CpuTask key={i} cpu={CPU} queues={newQueues} />;
            }
        })}
    </>);
}

function Err({ err }: { err: State<boolean> }) {
    if (err.Get()) {
        return <p className="err">Something went wrong</p>;
    }
    else {
        return <></>;
    }
}