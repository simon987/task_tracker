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
        return this.http.get(this.url + "/logs");
    }

    getProjectStats(id: number) {
        return this.http.get(this.url + "/project/stats/" + id)
    }
}
