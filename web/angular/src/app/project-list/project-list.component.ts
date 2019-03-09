import {Component, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {Project} from '../models/project';
import {AuthService} from '../auth.service';

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent implements OnInit {

    constructor(private apiService: ApiService,
                public authService: AuthService) {
    }

    projects: Project[];

    ngOnInit() {
        this.getProjects();
    }

    refresh() {
        this.getProjects();
    }

    getProjects() {
        this.apiService.getProjects().subscribe(data =>
            this.projects = data['content']['projects']);
    }

}
