import {Component, OnInit} from '@angular/core';
import {AuthService} from "../auth.service";

import * as moment from "moment"

@Component({
    selector: 'app-account-details',
    templateUrl: './account-details.component.html',
    styleUrls: ['./account-details.component.css']
})
export class AccountDetailsComponent implements OnInit {

    public moment = moment;

    constructor(public authService: AuthService) {
    }

    ngOnInit() {
    }

    public logout() {
        this.authService.logout();
    }
}
