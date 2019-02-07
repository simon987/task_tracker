import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute} from "@angular/router";


@Component({
    selector: 'app-project-dashboard',
    templateUrl: './project-dashboard.component.html',
    styleUrls: ['./project-dashboard.component.css']
})
export class ProjectDashboardComponent implements OnInit {

    private projectId;
    project: Project;


    constructor(private apiService: ApiService, private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
        })
    }

    private getProject() {
        this.apiService.getProject(this.projectId).subscribe(data => {
            this.project = <Project>{
                id: data["project"]["id"],
                name: data["project"]["name"],
                clone_url: data["project"]["clone_url"],
                git_repo: data["project"]["git_repo"],
                motd: data["project"]["motd"],
                priority: data["project"]["priority"],
                version: data["project"]["version"],
                public: data["project"]["public"],
            }
        })
    }
}
