import {Component, OnInit} from '@angular/core';
import {Project} from "../models/project";
import {ApiService} from "../api.service";
import {ActivatedRoute, Router} from "@angular/router";
import {MessengerService} from "../messenger.service";

@Component({
    selector: 'app-update-project',
    templateUrl: './update-project.component.html',
    styleUrls: ['./update-project.component.css']
})
export class UpdateProjectComponent implements OnInit {

    constructor(private apiService: ApiService,
                private route: ActivatedRoute,
                private messengerService: MessengerService,
                private router: Router) {
    }

    private project: Project;
    private projectId: number;

    ngOnInit() {
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

    onSubmit() {
        this.apiService.updateProject(this.project).subscribe(
            data => {
                this.router.navigateByUrl("/project/" + this.project.id);
            },
            error => {
                console.log(error.error.message);
                this.messengerService.show(error.error.message);
            }
        )
    }

}
