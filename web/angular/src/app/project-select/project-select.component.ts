import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {ApiService} from '../api.service';
import {Project} from '../models/project';

@Component({
    selector: 'project-select',
    templateUrl: './project-select.component.html',
    styleUrls: ['./project-select.component.css']
})
export class ProjectSelectComponent implements OnInit {

    projectList: Project[];

    @Input() project: Project;
    @Input() placeholder: string;
    @Output() projectChange = new EventEmitter<Project>();

    constructor(private apiService: ApiService) {
    }

    ngOnInit() {
    }

    loadProjectList() {
        this.apiService.getProjects().subscribe(data => {
            this.projectList = data['content']['projects'];
        });
    }
}
