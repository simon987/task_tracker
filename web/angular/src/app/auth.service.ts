import {Injectable} from '@angular/core';
import {ApiService} from "./api.service";
import {Credentials} from "./models/credentials";
import {MessengerService} from "./messenger.service";
import {Router} from "@angular/router";

@Injectable({
    providedIn: 'root'
})
export class AuthService {

    account: Manager;

    constructor(private apiService: ApiService,
                private messengerService: MessengerService,
                private router: Router) {
        this.apiService.getAccountDetails()
            .subscribe((data: any) => {
                this.account = data.manager;
            })
    }

    public login(credentials: Credentials) {
        return this.apiService.login(credentials)
            .subscribe(
                () => {
                    this.apiService.getAccountDetails()
                        .subscribe((data: any) => {
                            this.account = data.manager;
                            this.router.navigateByUrl("/account");
                        })
                },
                error => {
                    console.log(error);
                    this.messengerService.show(error.error.message);
                }
            )
    }

    public logout() {
        return this.apiService.logout()
            .subscribe(
                () => {
                    this.account = null;
                    this.router.navigateByUrl("");
                },
                error => {
                    console.log(error);
                    this.messengerService.show(error.error.message);
                }
            )
    }

    public register(credentials: Credentials) {
        return this.apiService.register(credentials)
            .subscribe(() =>
                    this.apiService.getAccountDetails()
                        .subscribe((data: any) => {
                            this.account = data.manager;
                            this.router.navigateByUrl("/account");
                        }),
                error => {
                    console.log(error);
                    this.messengerService.show(error.error.message);
                }
            )
    }
}
