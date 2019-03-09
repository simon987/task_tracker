import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ApiService} from '../api.service';
import {Manager} from '../models/manager';

@Component({
    selector: 'manager-select',
    templateUrl: './manager-select.component.html',
    styleUrls: ['./manager-select.component.css']
})
export class ManagerSelectComponent implements OnInit {

    manager: Manager;
    managerList: Manager[];

    @Output()
    managerChange = new EventEmitter<Manager>();

    constructor(private apiService: ApiService) {
    }

    ngOnInit() {
    }

    loadManagerList() {
        this.apiService.getManagerList()
            .subscribe(data => this.managerList = data['content']['managers']);
    }


}
