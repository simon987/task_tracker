import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {LogsComponent} from './logs/logs.component';

import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
    MatAutocompleteModule,
    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatDividerModule,
    MatExpansionModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatMenuModule,
    MatPaginatorIntl,
    MatPaginatorModule,
    MatSelectModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatSnackBarModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
    MatTreeModule
} from "@angular/material";
import {ApiService} from "./api.service";
import {MessengerService} from "./messenger.service";
import {HttpClient, HttpClientModule} from "@angular/common/http";
import {ProjectDashboardComponent} from './project-dashboard/project-dashboard.component';
import {ProjectListComponent} from './project-list/project-list.component';
import {CreateProjectComponent} from './create-project/create-project.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {UpdateProjectComponent} from './update-project/update-project.component';
import {SnackBarComponent} from "./messenger/snack-bar.component";
import {TranslateLoader, TranslateModule, TranslateService} from "@ngx-translate/core";
import {TranslateHttpLoader} from "@ngx-translate/http-loader";
import {TranslatedPaginator} from "./TranslatedPaginatorConfiguration";
import {LoginComponent} from './login/login.component';
import {CreateAccountComponent} from './create-account/create-account.component';


export function createTranslateLoader(http: HttpClient) {
    return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}


@NgModule({
    declarations: [
        AppComponent,
        LogsComponent,
        ProjectDashboardComponent,
        ProjectListComponent,
        CreateProjectComponent,
        UpdateProjectComponent,
        SnackBarComponent,
        LoginComponent,
        CreateAccountComponent,
    ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        MatMenuModule,
        MatIconModule,
        MatTableModule,
        MatPaginatorModule,
        MatSortModule,
        MatFormFieldModule,
        MatInputModule,
        MatToolbarModule,
        MatCardModule,
        MatButtonModule,
        MatAutocompleteModule,
        ReactiveFormsModule,
        FormsModule,
        MatExpansionModule,
        MatTreeModule,
        BrowserAnimationsModule,
        HttpClientModule,
        MatSliderModule,
        MatSlideToggleModule,
        MatCheckboxModule,
        MatDividerModule,
        MatSnackBarModule,
        TranslateModule.forRoot({
                loader: {
                    provide: TranslateLoader,
                    useFactory: (createTranslateLoader),
                    deps: [HttpClient]
                }
            }
        ),
        MatSelectModule

    ],
    exports: [],
    providers: [
        ApiService,
        MessengerService,
        {provide: MatPaginatorIntl, useFactory: TranslatedPaginator, deps: [TranslateService]}
    ],
    entryComponents: [
        SnackBarComponent,
    ],
    bootstrap: [AppComponent]
})
export class AppModule {
}
