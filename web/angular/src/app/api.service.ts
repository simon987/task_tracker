import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";

@Injectable()
export class ApiService {

    private url: string = "http://localhost:42901";

    constructor(
        private http: HttpClient,
    ) {
    }

    getLogs() {
        return this.http.post(this.url + "/logs", "{\"level\":\"info\", \"since\":10000}");
    }

    getProjectStats(id: number) {
        return this.http.get(this.url + "/project/stats/" + id)
    }

    getProjects() {
        return this.http.get(this.url + "/project/stats")
    }

    getProject(id: number) {
        return this.http.get(this.url + "/project/get/" + id)
    }
}
