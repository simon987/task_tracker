import {Worker} from './worker';
import {Project} from './project';

export class SubmitTaskOptions {
    public worker: Worker;
    public project: Project;
    public recipe: string;
    public uniqueStr: string;
    public maxAssignTime = 3600;
    public verificationCount = 0;
    public maxRetries = 3;
    public priority = 1;
}
