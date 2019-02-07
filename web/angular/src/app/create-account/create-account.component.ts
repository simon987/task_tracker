import {Component, OnInit} from '@angular/core';
import {Credentials} from "../models/credentials";

@Component({
    selector: 'app-create-account',
    templateUrl: './create-account.component.html',
    styleUrls: ['./create-account.component.css']
})
export class CreateAccountComponent implements OnInit {

    credentials: Credentials = <Credentials>{};

    constructor() {
    }

    ngOnInit() {
    }

    canCreate(): boolean {
        return this.credentials.username && this.credentials.username != "" &&
            this.credentials.password == this.credentials.repeatPassword
    }

    onClick() {
        alert("e")
    }

}
