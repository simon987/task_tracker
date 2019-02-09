import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {TranslateService} from "@ngx-translate/core";

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {

    langChange(lang: any) {
        this.translate.use(lang.lang)
    }

    langList: any[] = [
        {lang: "fr", display: "Français"},
        {lang: "en", display: "English"},
    ];

    constructor(private translate: TranslateService, private router: Router) {

        translate.addLangs([
            "en",
            "fr"
        ]);

        translate.setDefaultLang("en");
    }

}
