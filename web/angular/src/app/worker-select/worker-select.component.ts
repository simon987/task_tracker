import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {ApiService} from '../api.service';
import {Manager} from '../models/manager';
import {Worker} from '../models/worker';

@Component({
    selector: 'worker-select',
    templateUrl: './worker-select.component.html',
    styleUrls: ['./worker-select.component.css']
})
export class WorkerSelectComponent implements OnInit {

    @Input() worker: Worker;
    workerList: Worker[];

    @Output()
    workerChange = new EventEmitter<Worker>();

    constructor(private apiService: ApiService) {
    }

    ngOnInit() {
    }

    loadWorkerList() {
        this.apiService.getWorkerStats()
            .subscribe(data => {
                this.workerList = data['content']['stats'].sort((a, b) =>
                    (a.alias > b.alias) ? 1 : -1);
            });
    }
}
