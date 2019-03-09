import {Component, OnInit, ViewChild} from '@angular/core';
import {ApiService} from '../api.service';
import {MessengerService} from '../messenger.service';
import {TranslateService} from '@ngx-translate/core';
import {MatPaginator, MatSort, MatTableDataSource} from '@angular/material';

import * as moment from 'moment';
import {AuthService} from '../auth.service';
import {Manager} from '../models/manager';

@Component({
    selector: 'app-manager-list',
    templateUrl: './manager-list.component.html',
    styleUrls: ['./manager-list.component.css']
})
export class ManagerListComponent implements OnInit {

    managers = [];
    data;
    moment = moment;
    cols = ['username', 'tracker_admin', 'register_time', 'actions'];

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(private apiService: ApiService,
                private messengerService: MessengerService,
                private translate: TranslateService,
                private authService: AuthService
    ) {
        this.data = new MatTableDataSource<Manager>();
    }

    ngOnInit() {
        this.getManagers();
        this.data.paginator = this.paginator;
        this.data.sort = this.sort;
    }

    canPromote(manager: Manager) {
        return !manager.tracker_admin;
    }

    canDemote(manager: Manager) {
        return manager.tracker_admin && manager.username !== this.authService.account.username;
    }

    public promote(manager: Manager) {
        this.apiService.promote(manager.id)
            .subscribe(() => this.getManagers());
    }

    public demote(manager: Manager) {
        this.apiService.demote(manager.id)
            .subscribe(() => this.getManagers());
    }

    private getManagers() {
        this.apiService.getManagerList()
            .subscribe(data => {
                    this.data.data = data['content']['managers'];
                },
                error => {
                    if (error && (error.status === 401 || error.status === 403)) {
                        console.log(error.error.message);
                        this.translate.get('manager_list.unauthorized')
                            .subscribe(t => this.messengerService.show(t));
                    }
                });
    }
}
