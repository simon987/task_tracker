import {Worker} from "../models/worker"
import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute, Router} from "@angular/router";
import {MessengerService} from "../messenger.service";
import {TranslateService} from "@ngx-translate/core";

import * as moment from "moment"

@Component({
    selector: 'app-project-perms',
    templateUrl: './project-perms.component.html',
    styleUrls: ['./project-perms.component.css']
})
export class ProjectPermsComponent implements OnInit {

    constructor(private apiService: ApiService,
                private route: ActivatedRoute,
                private messengerService: MessengerService,
                private translate: TranslateService,
                private router: Router) {
    }

    project: Project;
    private projectId: number;
    requests: Worker[];
    unauthorized: boolean = false;
    moment = moment;

    ngOnInit() {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
            this.getProjectRequests();
        })
    }

    private getProject() {
        this.apiService.getProject(this.projectId).subscribe(data => {
            this.project = data["project"]
        })
    }

    private getProjectRequests() {
        this.apiService.getProjectAccessRequests(this.projectId).subscribe(
            data => {
                this.requests = data["requests"]
            },
            error => {
                if (error && (error.status == 401 || error.status == 403)) {
                    this.unauthorized = true;
                }
            })
    }

    private refresh() {
        this.getProjectRequests()
    }
}
