import {NgModule} from '@angular/core';
import {NavigationEnd, NavigationStart, Router, RouterModule, Routes} from '@angular/router';
import {LogsComponent} from "./logs/logs.component";
import {ProjectDashboardComponent} from "./project-dashboard/project-dashboard.component";
import {ProjectListComponent} from "./project-list/project-list.component";
import {CreateProjectComponent} from "./create-project/create-project.component";
import {UpdateProjectComponent} from "./update-project/update-project.component";
import {Title} from "@angular/platform-browser";
import {filter} from "rxjs/operators";
import {TranslateService} from "@ngx-translate/core";

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

    constructor(private title: Title, private router: Router, private translate: TranslateService) {
        router.events
            .pipe(filter(event => event instanceof NavigationEnd))
            .subscribe((event: NavigationStart) => {
                    this.updateTitle(translate, title, event.url)
                }
            );

        translate.onLangChange.subscribe(() =>
            this.updateTitle(translate, title, router.url)
        )
    }

    private updateTitle(tr: TranslateService, title: Title, url: string) {
        url = url.substr(1);
        tr.get("title." + url.substring(0, url.indexOf("/") == -1 ? url.length : url.indexOf("/")))
            .subscribe((t) => title.setTitle(t))
    }
}

