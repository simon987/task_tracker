import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {LogsComponent} from "./logs/logs.component";
import {ProjectDashboardComponent} from "./project-dashboard/project-dashboard.component";

const routes: Routes = [
    {path: "log", component: LogsComponent},
    {path: "project", component: ProjectDashboardComponent}
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
