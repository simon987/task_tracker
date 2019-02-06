import {Component} from '@angular/core';
import {TranslateService} from "@ngx-translate/core";
import {MatSelectChange} from "@angular/material";

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {

    langChange(event: MatSelectChange) {
        this.translate.use(event.value)
    }

    langList: any[] = [
        {lang: "fr", display: "Français"},
        {lang: "en", display: "English"},
    ];

    constructor(private translate: TranslateService) {

        translate.addLangs([
            "en",
            "fr"
        ]);

        translate.setDefaultLang("en");
    }

}
