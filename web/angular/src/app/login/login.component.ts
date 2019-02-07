import {Component, OnInit} from '@angular/core';
import {Credentials} from "../models/credentials";

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    credentials: Credentials = <Credentials>{};

    constructor() {
    }

    ngOnInit() {
    }

    onClick() {

    }
}
