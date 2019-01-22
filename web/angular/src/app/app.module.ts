import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {LogsComponent} from './logs/logs.component';

import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
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

@NgModule({
    declarations: [
        AppComponent,
        LogsComponent,
        ProjectDashboardComponent
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
