import {Component, OnInit} from '@angular/core';
import {AuthService} from "../auth.service";

@Component({
    selector: 'app-account-details',
    templateUrl: './account-details.component.html',
    styleUrls: ['./account-details.component.css']
})
export class AccountDetailsComponent implements OnInit {

    constructor(private authService: AuthService) {
    }

    ngOnInit() {
    }

}
