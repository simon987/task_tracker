import {Component, OnInit} from '@angular/core';
import {Credentials} from "../models/credentials";
import {ApiService} from "../api.service";
import {MessengerService} from "../messenger.service";
import {Router} from "@angular/router";
import {AuthService} from "../auth.service";

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    credentials: Credentials = <Credentials>{};

    constructor(private apiService: ApiService,
                private messengerService: MessengerService,
                private router: Router,
                private authService: AuthService) {
    }

    ngOnInit() {
    }

    login() {
        this.authService.login(this.credentials)
    }

    register() {
        this.apiService.register(this.credentials)
            .subscribe(
                () => {
                    this.router.navigateByUrl("/account")
                },
                error => {
                    console.log(error);
                    this.messengerService.show(error.error.message);
                }
            )
    }

    canCreate(): boolean {
        return this.credentials.username && this.credentials.username != "" &&
            this.credentials.password == this.credentials.repeatPassword
    }
}
