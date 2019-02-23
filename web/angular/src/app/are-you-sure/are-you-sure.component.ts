import {Component, OnInit} from '@angular/core';
import {MatDialogRef} from "@angular/material";

@Component({
    selector: 'app-are-you-sure',
    templateUrl: './are-you-sure.component.html',
    styleUrls: ['./are-you-sure.component.css']
})
export class AreYouSureComponent implements OnInit {

    constructor(public dialogRef: MatDialogRef<AreYouSureComponent>) {
    }

    ngOnInit() {

    }

    onNoClick() {
        this.dialogRef.close(false)
    }

    onYesClick() {
        this.dialogRef.close(true)
    }
}
