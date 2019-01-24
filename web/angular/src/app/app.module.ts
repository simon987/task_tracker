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
    MatExpansionModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatMenuModule,
    MatPaginatorModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
    MatTreeModule
} from "@angular/material";
import {ApiService} from "./api.service";
import {HttpClientModule} from "@angular/common/http";
import {ProjectDashboardComponent} from './project-dashboard/project-dashboard.component';
import {ProjectListComponent} from './project-list/project-list.component';
import {CreateProjectComponent} from './create-project/create-project.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {UpdateProjectComponent} from './update-project/update-project.component';

@NgModule({
    declarations: [
        AppComponent,
        LogsComponent,
        ProjectDashboardComponent,
        ProjectListComponent,
        CreateProjectComponent,
        UpdateProjectComponent
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
    ],
    exports: [],
    providers: [
        ApiService,
    ],
    bootstrap: [AppComponent]
})
export class AppModule {
}
