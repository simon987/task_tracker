import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent implements OnInit {

    constructor(private apiService: ApiService) {
    }

    projects: Project[];

    ngOnInit() {
        this.getProjects()
    }

    getProjects() {
        this.apiService.getProjects().subscribe(data => this.projects = data["projects"]);
    }

}
