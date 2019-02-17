import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute} from "@angular/router";

import * as moment from "moment"
import {WorkerAccess} from "../models/worker-access";
import {AuthService} from "../auth.service";
import {Manager, ManagerRoleOnProject} from "../models/manager";
import {MessengerService} from "../messenger.service";
import {TranslateService} from "@ngx-translate/core";

@Component({
    selector: 'app-project-perms',
    templateUrl: './project-perms.component.html',
    styleUrls: ['./project-perms.component.css']
})
export class ProjectPermsComponent implements OnInit {

    constructor(private apiService: ApiService,
                private route: ActivatedRoute,
                private translate: TranslateService,
                private messenger: MessengerService,
                public auth: AuthService) {
    }

    project: Project;
    private projectId: number;
    accesses: WorkerAccess[];
    managerRoles: ManagerRoleOnProject;
    unauthorized: boolean = false;
    moment = moment;

    ngOnInit() {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
            this.getProjectAccesses();
            this.getProjectManagers();
        })
    }

    public acceptRequest(wa: WorkerAccess) {
        this.apiService.acceptWorkerAccessRequest(wa.worker.id, this.projectId)
            .subscribe(() => {
                this.getProjectAccesses();
                this.translate.get("perms.set").subscribe(t => this.messenger.show(t));
            })
    }

    public rejectRequest(wa: WorkerAccess) {
        this.apiService.rejectWorkerAccessRequest(wa.worker.id, this.projectId)
            .subscribe(() => {
                this.getProjectAccesses();
                this.translate.get("perms.set").subscribe(t => this.messenger.show(t));
            })
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

    private getProjectManagers() {
        this.apiService.getManagerListWithRoleOn(this.projectId)
            .subscribe(data => {
                this.managerRoles = data["content"]["managers"].map(d =>
                    ManagerRoleOnProject.fromEntity(d))
            })
    }

    public refresh() {
        this.getProjectAccesses();
        this.getProjectManagers();
    }

    public onSelectManager(manager: Manager) {
        if (manager.id != this.auth.account.id) {
            this.apiService.setManagerRoleOnProject(this.projectId, 1, manager.id)
                .subscribe(() => this.refresh())
        }
    }

    public onRoleChange(manager: ManagerRoleOnProject) {
        this.apiService.setManagerRoleOnProject(this.projectId, manager.role, manager.manager.id)
            .subscribe(() => {
                this.refresh();
                this.translate.get("perms.set").subscribe(t => this.messenger.show(t));
            })
    }
}
