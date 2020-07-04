import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Project} from './models/project';
import {Worker} from './models/worker';
import {Credentials} from './models/credentials';
import {SubmitTaskOptions} from './models/console';

@Injectable()
export class ApiService {

    public url: string = window.location.protocol + '//' + window.location.host + '/api';
    private options: {
        withCredentials: true,
        responseType: 'json'
    };

    constructor(
        private http: HttpClient,
    ) {
    }

    private static getWorkerHeaders(w: Worker): HttpHeaders {
        return new HttpHeaders({
            'X-Worker-ID': w.id.toString(),
            'X-Secret': w.secret,
        });
    }

    getLogs(level: number) {
        return this.http.post(this.url + '/logs', {level: level, since: 1}, this.options);
    }

    getProjects() {
        return this.http.get(this.url + '/project/list', this.options);
    }

    getProject(id: number) {
        return this.http.get(this.url + '/project/get/' + id, this.options);
    }

    createProject(project: Project) {
        return this.http.post(this.url + '/project/create', project, this.options);
    }

    updateProject(project: Project) {
        return this.http.post(this.url + '/project/update/' + project.id, project, this.options);
    }

    register(credentials: Credentials) {
        return this.http.post(this.url + '/register', credentials, this.options);
    }

    login(credentials: Credentials) {
        return this.http.post(this.url + '/login', credentials, this.options);
    }

    logout() {
        return this.http.get(this.url + '/logout', this.options);
    }

    getAccountDetails() {
        return this.http.get(this.url + '/account', this.options);
    }

    getMonitoringSnapshots(count: number, project: number) {
        return this.http.get(this.url + `/project/monitoring/${project}?count=${count}`, this.options);
    }

    getAssigneeStats(project: number) {
        return this.http.get(this.url + `/project/assignees/${project}`, this.options);
    }

    getWorkerStats() {
        return this.http.get(this.url + `/worker/stats`, this.options);
    }

    getProjectAccess(project: number) {
        return this.http.get(this.url + `/project/access_list/${project}`, this.options);
    }

    getManagerList() {
        return this.http.get(this.url + '/manager/list', this.options);
    }

    getManagerListWithRoleOn(project: number) {
        return this.http.get(this.url + '/manager/list_for_project/' + project, this.options);
    }

    promote(managerId: number) {
        return this.http.get(this.url + `/manager/promote/${managerId}`, this.options);
    }

    demote(managerId: number) {
        return this.http.get(this.url + `/manager/demote/${managerId}`, this.options);
    }

    acceptWorkerAccessRequest(wid: number, pid: number) {
        return this.http.post(this.url + `/project/accept_request/${pid}/${wid}`, null, this.options);
    }

    rejectWorkerAccessRequest(wid: number, pid: number) {
        return this.http.post(this.url + `/project/reject_request/${pid}/${wid}`, null, this.options);
    }

    setManagerRoleOnProject(pid: number, role: number, manager: number) {
        return this.http.post(this.url + `/manager/set_role_for_project/${pid}`,
            {'role': role, 'manager': manager}, this.options);
    }

    getSecret(pid: number) {
        return this.http.get(this.url + `/project/secret/${pid}`, this.options);
    }

    setSecret(pid: number, secret: string) {
        return this.http.post(this.url + `/project/secret/${pid}`, {'secret': secret}, this.options);
    }

    getWebhookSecret(pid: number) {
        return this.http.get(this.url + `/project/webhook_secret/${pid}`, this.options);
    }

    setWebhookSecret(pid: number, secret: string) {
        return this.http.post(this.url + `/project/webhook_secret/${pid}`, {'webhook_secret': secret}, this.options);
    }

    resetFailedTasks(pid: number) {
        return this.http.post(this.url + `/project/reset_failed_tasks/${pid}`, null, this.options);
    }

    hardReset(pid: number) {
        return this.http.post(this.url + `/project/hard_reset/${pid}`, null, this.options);
    }

    reclaimAssignedTasks(pid: number) {
        return this.http.post(this.url + `/project/reclaim_assigned_tasks/${pid}`, null, this.options);
    }

    workerSetPaused(wid: number, paused: boolean) {
        return this.http.post(this.url + '/worker/set_paused',
            {'worker': wid, 'paused': paused}, this.options);
    }

    getWorker(wid: number) {
        return this.http.get(this.url + `/worker/get/${wid}`, this.options);
    }

    workerSubmitTask(taskOptions: SubmitTaskOptions) {
        return this.http.post(this.url + `/task/submit`, {
            project: taskOptions.project.id,
            max_retries: taskOptions.maxRetries,
            recipe: taskOptions.recipe,
            priority: taskOptions.priority,
            max_assign_time: taskOptions.maxAssignTime,
            hash64: 0,
            unique_string: taskOptions.uniqueStr,
            verification_count: taskOptions.verificationCount
        }, {
            headers: ApiService.getWorkerHeaders(taskOptions.worker),
            responseType: 'json'
        });
    }
}
