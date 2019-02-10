import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Project} from "./models/project";
import {Credentials} from "./models/credentials";

@Injectable()
export class ApiService {

    private url: string = "http://localhost/api";
    private options: {
        withCredentials: true,
        responseType: "json"
    };

    constructor(
        private http: HttpClient,
    ) {
    }

    getLogs(level: number) {
        return this.http.post(this.url + "/logs", {level: level, since: 1}, this.options);
    }

    getProjects() {
        return this.http.get(this.url + "/project/list", this.options)
    }

    getProject(id: number) {
        return this.http.get(this.url + "/project/get/" + id, this.options)
    }

    createProject(project: Project) {
        return this.http.post(this.url + "/project/create", project, this.options)
    }

    updateProject(project: Project) {
        return this.http.post(this.url + "/project/update/" + project.id, project, this.options)
    }

    register(credentials: Credentials) {
        return this.http.post(this.url + "/register", credentials, this.options)
    }

    login(credentials: Credentials) {
        return this.http.post(this.url + "/login", credentials, this.options)
    }

    logout() {
        return this.http.get(this.url + "/logout", this.options)
    }

    getAccountDetails() {
        return this.http.get(this.url + "/account", this.options)
    }

    getMonitoringSnapshots(count: number, project: number) {
        return this.http.get(this.url + `/project/monitoring/${project}?count=${count}`, this.options)
    }

    getAssigneeStats(project: number) {
        return this.http.get(this.url + `/project/assignees/${project}`, this.options)
    }

    getWorkerStats() {
        return this.http.get(this.url + `/worker/stats`, this.options)
    }

}
