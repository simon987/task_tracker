import {Component, OnInit} from '@angular/core';
import {Project} from "../models/project";

@Component({
    selector: 'app-create-project',
    templateUrl: './create-project.component.html',
    styleUrls: ['./create-project.component.css']
})
export class CreateProjectComponent implements OnInit {

    private project = new Project();

    constructor() {
        this.project.name = "test";
        this.project.public = true;
    }

    ngOnInit() {
    }

}
