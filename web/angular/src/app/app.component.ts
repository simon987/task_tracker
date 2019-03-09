import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {TranslateService} from '@ngx-translate/core';
import {AuthService} from './auth.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {

    constructor(private translate: TranslateService,
                public router: Router,
                public authService: AuthService) {

        translate.addLangs([
            'en',
            'fr'
        ]);

        translate.setDefaultLang('en');
    }

    langList: any[] = [
        {lang: 'fr', display: 'Français'},
        {lang: 'en', display: 'English'},
    ];

    langChange(lang: any) {
        this.translate.use(lang.lang);
    }

}
