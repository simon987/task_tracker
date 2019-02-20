import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {LogsComponent} from './logs/logs.component';

import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
    MatAutocompleteModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
    MatCheckboxModule,
    MatDividerModule,
    MatExpansionModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatPaginatorIntl,
    MatPaginatorModule,
    MatProgressBarModule,
    MatSelectModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatSnackBarModule,
    MatSortModule,
    MatStepperModule,
    MatTableModule,
    MatTabsModule,
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
import {AccountDetailsComponent} from './account-details/account-details.component';
import {WorkerDashboardComponent} from './worker-dashboard/worker-dashboard.component';
import {ProjectPermsComponent} from './project-perms/project-perms.component';
import {ManagerListComponent} from './manager-list/manager-list.component';
import {ProjectSelectComponent} from './project-select/project-select.component';
import {ManagerSelectComponent} from './manager-select/manager-select.component';
import {ProjectIconComponent} from './project-icon/project-icon.component';
import {IndexComponent} from './index/index.component';
import {ProjectSecretComponent} from './project-secret/project-secret.component';


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
        AccountDetailsComponent,
        WorkerDashboardComponent,
        ProjectPermsComponent,
        ManagerListComponent,
        ProjectSelectComponent,
        ManagerSelectComponent,
        ProjectIconComponent,
        IndexComponent,
        ProjectSecretComponent,
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
        MatSelectModule,
        MatProgressBarModule,
        MatTabsModule,
        MatListModule,
        MatButtonToggleModule,
        MatStepperModule

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
