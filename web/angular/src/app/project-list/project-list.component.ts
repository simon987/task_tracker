import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent implements OnInit {

    constructor(private apiService: ApiService) {
    }

    projects: any[];

    ngOnInit() {
        this.getProjects()
    }

    getProjects() {
        this.apiService.getProjects().subscribe(data => this.projects = data["stats"]);
    }

}
