import {Component, OnInit} from '@angular/core';
import {Project} from "../models/project";
import {ApiService} from "../api.service";
import {MessengerService} from "../messenger.service";
import {Router} from "@angular/router";


@Component({
    selector: 'app-create-project',
    templateUrl: './create-project.component.html',
    styleUrls: ['./create-project.component.css']
})
export class CreateProjectComponent implements OnInit {

    project = new Project();

    constructor(private apiService: ApiService,
                private messengerService: MessengerService,
                private router: Router) {
        this.project.name = "test";
        this.project.public = true;
    }

    ngOnInit() {
    }

    onSubmit() {
        this.apiService.createProject(this.project).subscribe(
            data => {
                this.router.navigateByUrl("/project/" + data["id"]);
            },
            error => {
                console.log(error.error.message);
                this.messengerService.show(error.error.message);
            }
        )
    }

}
