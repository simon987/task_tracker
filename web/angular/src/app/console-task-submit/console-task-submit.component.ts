import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {SubmitTaskOptions} from '../models/console';
import {ApiService} from '../api.service';

@Component({
    selector: 'app-console-task-submit',
    templateUrl: './console-task-submit.component.html',
    styleUrls: ['./console-task-submit.component.css']
})
export class ConsoleTaskSubmitComponent implements OnInit {

    public submitOptions = new SubmitTaskOptions();

    @Output() submitTask = new EventEmitter<SubmitTaskOptions>();

    constructor(private apiService: ApiService) {
    }

    ngOnInit() {
    }

    onSubmit() {
        this.apiService.getWorker(this.submitOptions.worker.id).subscribe(data => {
            this.submitOptions.worker.secret = data['content']['worker']['secret'];
            this.submitTask.emit(this.submitOptions);
        });
    }

    public buttonDisabled(): boolean {
        return this.submitOptions.project === undefined || this.submitOptions.worker === undefined;
    }
}
