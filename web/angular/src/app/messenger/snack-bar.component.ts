import {Component, OnInit} from '@angular/core';
import {MessengerService} from "../messenger.service";
import {MessengerState} from "./messenger";
import {Subscription} from "rxjs";
import {MatSnackBar, MatSnackBarConfig} from "@angular/material";

@Component({
    selector: 'messenger-snack-bar',
    templateUrl: 'messenger-snack-bar.html',
    styleUrls: ['messenger-snack-bar.css'],
})
export class SnackBarComponent implements OnInit {

    private subscription: Subscription;

    constructor(private messengerService: MessengerService, private snackBar: MatSnackBar) {

    }

    ngOnInit() {
        this.subscription = this.messengerService.messengerSubject
            .subscribe((state: MessengerState) => {
                if (state.hidden) {
                    this.snackBar.dismiss();
                } else {
                    this.snackBar.open(state.message, "Close", <MatSnackBarConfig>{
                        duration: 10 * 1000,
                    })
                }
            });
    }
}
