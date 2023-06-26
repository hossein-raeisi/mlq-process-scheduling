import Queue from "./Queue";
import { CPUUpdate, AddProcess } from "./Updates";
export function CpuTask({ cpu, queues }: { cpu: CPUUpdate, queues: AddProcess[][] }) {
    const a = 1;
    return (
        <div className={"task "}>
            <CPUUsage cpu={cpu} />
            {queues.map((queue, i) => {
                return <Queue key={i} i={i} queue={queue} />;
            })}
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
                    <td>{cpu.StartStr}</td>
                </tr>
                <tr>
                    <td>process</td>
                    <td>{cpu.Name}</td>
                </tr>
                <tr>
                    <td>end</td>
                    <td>{cpu.EndStr}</td>
                </tr>
            </tbody>
        </table>
    );
}