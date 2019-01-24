import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {LogsComponent} from "./logs/logs.component";
import {ProjectDashboardComponent} from "./project-dashboard/project-dashboard.component";
import {ProjectListComponent} from "./project-list/project-list.component";
import {CreateProjectComponent} from "./create-project/create-project.component";
import {UpdateProjectComponent} from "./update-project/update-project.component";

const routes: Routes = [
    {path: "log", component: LogsComponent},
    {path: "projects", component: ProjectListComponent},
    {path: "project/:id", component: ProjectDashboardComponent},
    {path: "project/:id/update", component: UpdateProjectComponent},
    {path: "new_project", component: CreateProjectComponent}
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
