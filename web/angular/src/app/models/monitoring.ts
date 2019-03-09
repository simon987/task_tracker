export interface MonitoringSnapshot {
    new_task_count: number;
    failed_task_count: number;
    closed_task_count: number;
    awaiting_verification_count: number;
    time_stamp: number;
}

export interface AssignedTasks {
    assignee: string;
    task_count: number;
}
