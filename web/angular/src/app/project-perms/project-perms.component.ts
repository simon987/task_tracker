import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute, Router} from "@angular/router";
import {MessengerService} from "../messenger.service";
import {TranslateService} from "@ngx-translate/core";

import * as moment from "moment"
import {WorkerAccess} from "../models/worker-access";

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
    accesses: WorkerAccess[];
    unauthorized: boolean = false;
    moment = moment;

    ngOnInit() {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
            this.getProjectAccesses();
        })
    }

    public acceptRequest(wa: WorkerAccess) {
        this.apiService.acceptWorkerAccessRequest(wa.worker.id, this.projectId)
            .subscribe(() => this.getProjectAccesses())
    }

    public rejectRequest(wa: WorkerAccess) {
        this.apiService.rejectWorkerAccessRequest(wa.worker.id, this.projectId)
            .subscribe(() => this.getProjectAccesses())
    }

    private getProject() {
        this.apiService.getProject(this.projectId).subscribe(data => {
            this.project = data["content"]["project"]
        })
    }

    private getProjectAccesses() {
        this.apiService.getProjectAccess(this.projectId).subscribe(
            data => {
                this.accesses = data["content"]["accesses"]
            },
            error => {
                if (error && (error.status == 401 || error.status == 403)) {
                    this.unauthorized = true;
                }
            })
    }

    public refresh() {
        this.getProjectAccesses()
    }
}
