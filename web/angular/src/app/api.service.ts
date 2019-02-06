import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Project} from "./models/project";

@Injectable()
export class ApiService {

    private url: string = "http://localhost:42901";

    constructor(
        private http: HttpClient,
    ) {
    }

    getLogs() {
        return this.http.post(this.url + "/logs", "{\"level\":6, \"since\":1}");
    }

    getProjects() {
        return this.http.get(this.url + "/project/list")
    }

    getProject(id: number) {
        return this.http.get(this.url + "/project/get/" + id)
    }

    createProject(project: Project) {
        return this.http.post(this.url + "/project/create", project)
    }

    updateProject(project: Project) {
        return this.http.post(this.url + "/project/update/" + project.id, project)
    }
}
