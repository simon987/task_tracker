import {Worker} from './worker';

export interface WorkerAccess {
    submit: boolean;
    assign: boolean;
    request: boolean;
    worker: Worker;
    project: number;
}
