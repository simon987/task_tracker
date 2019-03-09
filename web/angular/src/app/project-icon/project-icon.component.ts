import {Component, Input, OnInit} from '@angular/core';
import {Project} from '../models/project';

@Component({
    selector: 'project-icon',
    templateUrl: './project-icon.component.html',
    styleUrls: ['./project-icon.component.css']
})
export class ProjectIconComponent implements OnInit {


    @Input()
    project: Project;

    constructor() {
    }

    ngOnInit() {
    }

}
