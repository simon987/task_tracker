import {Component, OnInit} from '@angular/core';
import {Project} from '../models/project';
import {ApiService} from '../api.service';
import {MessengerService} from '../messenger.service';
import {Router} from '@angular/router';
import {AuthService} from '../auth.service';


@Component({
    selector: 'app-create-project',
    templateUrl: './create-project.component.html',
    styleUrls: ['./create-project.component.css']
})
export class CreateProjectComponent implements OnInit {

    project = <Project>{};
    selectedProject: Project = null;

    constructor(private apiService: ApiService,
                private messengerService: MessengerService,
                public authService: AuthService,
                private router: Router) {
    }

    ngOnInit() {
    }

    cloneUrlChange() {
        const tokens = this.project.clone_url.split('/');

        if (tokens.length > 2) {
            this.project.git_repo = tokens[tokens.length - 2] + '/' + tokens[tokens.length - 1];
        }
    }

    onSubmit() {
        this.project.chain = this.selectedProject ? this.selectedProject.id : 0;

        this.apiService.createProject(this.project).subscribe(
            data => {
                this.router.navigateByUrl('/project/' + data['content']['id']);
            },
            error => {
                console.log(error.error.message);
                this.messengerService.show(error.error.message);
            }
        );
    }

}
