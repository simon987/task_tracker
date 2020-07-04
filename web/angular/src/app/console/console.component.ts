import {Component, OnInit} from '@angular/core';
import {SubmitTaskOptions} from '../models/console';
import {ApiService} from '../api.service';
import {MessengerService} from '../messenger.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-console',
    templateUrl: './console.component.html',
    styleUrls: ['./console.component.css']
})
export class ConsoleComponent implements OnInit {

    constructor(private apiService: ApiService,
                private messenger: MessengerService,
                private translate: TranslateService) {
    }

    ngOnInit() {
    }

    onTaskSubmit(options: SubmitTaskOptions) {
        this.apiService.workerSubmitTask(options).subscribe(data => {
            this.translate.get('console.submit_ok').subscribe(t => this.messenger.show(t));
        }, error => {
            this.messenger.show(error.error.message);
        });
    }
}
