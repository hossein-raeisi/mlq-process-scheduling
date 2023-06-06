import Queue from "./Queue";
import { CPUUpdate, AddProcess } from "./Updates";
export default function Task({ cpu, queues }: { cpu: CPUUpdate, queues: AddProcess[][] }) {
    return (
        <div className={"task "} >
            <CPUUsage cpu={cpu} />
            {queues.map((queue, i) => <Queue key={i} i={i} queue={queue} />)}
        </div>
    );
}

function CPUUsage({cpu }: {cpu: CPUUpdate}) {
    return (
        <table>
            <thead>
                <tr>
                    <th>CPU Usage</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>start</td>
                    <td>{cpu.Start}</td>
                </tr>
                <tr>
                    <td>process</td>
                    <td>{cpu.Name}</td>
                </tr>
                <tr>
                    <td>end</td>
                    <td>{cpu.End}</td>
                </tr>
            </tbody>
        </table>
    );
}