import { AddProcess } from "./Updates";

export default function Queue({ queue, i }: { queue: AddProcess[], i: number }) {
    return (<table>
        <thead>
            <tr>
                <th colSpan={2}>{"Queue " + i}</th>
            </tr>
            <tr>
                <th>Name</th>
                <th>CBT</th>
            </tr>
        </thead>
        <tbody>
            {queue.map((process, i) => <Process key={i} process={process} />)}
        </tbody>
    </table>);
}

function Process({ process }: { process: AddProcess }) {
    return (<tr>
        <td>{process.Name}</td>
        <td>{process.CBT}</td>
    </tr>);
}