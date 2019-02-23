import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute} from "@angular/router";

import {Chart} from "chart.js";
import {AssignedTasks, MonitoringSnapshot} from "../models/monitoring";
import {TranslateService} from "@ngx-translate/core";
import {MessengerService} from "../messenger.service";
import {AuthService} from "../auth.service";
import {MatDialog} from "@angular/material";
import {AreYouSureComponent} from "../are-you-sure/are-you-sure.component";


@Component({
    selector: 'app-project-dashboard',
    templateUrl: './project-dashboard.component.html',
    styleUrls: ['./project-dashboard.component.css']
})
export class ProjectDashboardComponent implements OnInit {

    private projectId;
    project: Project;
    noTasks = false;

    private timeline: Chart;
    private statusPie: Chart;
    private assigneesPie: Chart;


    private colors = {
        new: "#76FF03",
        failed: "#FF3D00",
        closed: "#E0E0E0",
        awaiting: "#FFB74D",
        random: [
            "#3D5AFE", "#2979FF", "#2196F3",
            "#7C4DFF", "#673AB7", "#7C4DFF",
            "#FFC400", "#FFD740", "#FFC107",
            "#FF3D00", "#FF6E40", "#FF5722",
            "#76FF03", "#B2FF59", "#8BC34A"
        ]
    };

    snapshots: MonitoringSnapshot[] = [];
    lastSnapshot: MonitoringSnapshot;
    assignees: AssignedTasks[];

    constructor(private apiService: ApiService,
                private route: ActivatedRoute,
                private translate: TranslateService,
                public auth: AuthService,
                public dialog: MatDialog,
                private messenger: MessengerService) {
    }

    ngOnInit(): void {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
        });
    }

    public isSafeUrl(url: string) {
        if (url.substr(0, "http".length) == "http") {
            return true
        }
    }

    public refresh() {

        this.apiService.getMonitoringSnapshots(60, this.projectId)
            .subscribe((data: any) => {
                this.snapshots = data.content.snapshots;
                this.lastSnapshot = this.snapshots ? this.snapshots.sort((a, b) => {
                    return b.time_stamp - a.time_stamp
                })[0] : null;

                if (this.lastSnapshot == null || (this.lastSnapshot.awaiting_verification_count == 0 &&
                    this.lastSnapshot.closed_task_count == 0 &&
                    this.lastSnapshot.new_task_count == 0 &&
                    this.lastSnapshot.failed_task_count == 0)) {
                    this.noTasks = true;
                    return
                }
                this.noTasks = false;

                this.timeline.data.labels = this.snapshots.map(s => s.time_stamp * 1000 as any);
                this.timeline.data.datasets = this.makeTimelineDataset(this.snapshots);
                this.timeline.update();
                this.statusPie.data.datasets = [
                    {
                        label: "Task status",
                        data: [
                            this.lastSnapshot.new_task_count,
                            this.lastSnapshot.failed_task_count,
                            this.lastSnapshot.closed_task_count,
                            this.lastSnapshot.awaiting_verification_count,
                        ],
                        backgroundColor: [
                            this.colors.new,
                            this.colors.failed,
                            this.colors.closed,
                            this.colors.awaiting
                        ],
                    }
                ];
                this.statusPie.update();

                this.apiService.getAssigneeStats(this.projectId)
                    .subscribe((data: any) => {
                        this.assignees = data.content.assignees;
                        let colors = this.assignees.map(() => {
                            return this.colors.random[Math.floor(Math.random() * this.colors.random.length)]
                        });
                        this.assigneesPie.data.labels = this.assignees.map(x => x.assignee);
                        this.assigneesPie.data.datasets = [
                            {
                                label: "Task status",
                                data: this.assignees.map(x => x.task_count),
                                backgroundColor: colors,
                            }
                        ];
                        this.assigneesPie.update();
                    });
            })
    }

    private makeTimelineDataset(snapshots: MonitoringSnapshot[]) {
        return [
            {
                label: "New",
                type: "line",
                fill: false,
                borderColor: this.colors.new,
                backgroundColor: this.colors.new,
                data: snapshots.map(s => s.new_task_count),
                pointRadius: 0,
                lineTension: 0.2,
            },
            {
                label: "Failed",
                type: "line",
                fill: false,
                borderColor: this.colors.failed,
                backgroundColor: this.colors.failed,
                data: snapshots.map(s => s.failed_task_count),
                pointRadius: 0,
                lineTension: 0.2,
            },
            {
                label: "Closed",
                type: "line",
                fill: false,
                borderColor: this.colors.closed,
                backgroundColor: this.colors.closed,
                pointRadius: 0,
                data: snapshots.map(s => s.closed_task_count),
                lineTension: 0.2,
            },
            {
                label: "Awaiting verification",
                type: "line",
                fill: false,
                borderColor: this.colors.awaiting,
                backgroundColor: this.colors.awaiting,
                data: snapshots.map(s => s.awaiting_verification_count),
                pointRadius: 0,
                lineTension: 0.2,
            },
        ]
    }

    private setupTimeline() {
        let elem = document.getElementById("timeline") as any;
        let ctx = elem.getContext("2d");

        this.timeline = new Chart(ctx, {
            type: "bar",
            data: {
                labels: this.snapshots.map(s => s.time_stamp * 1000 as any),
                datasets: this.makeTimelineDataset(this.snapshots),
            },
            options: {
                title: {
                    display: true,
                    text: "Task status timeline",
                    position: "bottom"
                },
                legend: {
                    position: 'left',
                },
                scales: {
                    xAxes: [{
                        type: "time",
                        distribution: "series",
                        ticks: {
                            source: "auto"
                        },
                    }]
                },
                tooltips: {
                    enabled: true,
                    intersect: false,
                    mode: "index",
                    position: "nearest",
                },
                responsive: true
            }
        })
    }

    private setupStatusPie() {

        if (this.lastSnapshot == undefined || (this.lastSnapshot.awaiting_verification_count == 0 &&
            this.lastSnapshot.closed_task_count == 0 &&
            this.lastSnapshot.new_task_count == 0 &&
            this.lastSnapshot.failed_task_count == 0)) {
            this.noTasks = true;

            this.lastSnapshot = {
                closed_task_count: 0, time_stamp: 0, failed_task_count: 0,
                new_task_count: 0, awaiting_verification_count: 0
            }
        }

        let elem = document.getElementById("status-pie") as any;
        let ctx = elem.getContext("2d");

        this.statusPie = new Chart(ctx, {
            type: "doughnut",
            data: {
                labels: [
                    "New",
                    "Failed",
                    "Closed",
                    "Awaiting verification",
                ],
                datasets: [
                    {
                        label: "Task status",
                        data: [
                            this.lastSnapshot.new_task_count,
                            this.lastSnapshot.failed_task_count,
                            this.lastSnapshot.closed_task_count,
                            this.lastSnapshot.awaiting_verification_count,
                        ],
                        backgroundColor: [
                            this.colors.new,
                            this.colors.failed,
                            this.colors.closed,
                            this.colors.awaiting
                        ],
                    }
                ],
            },
            options: {
                responsive: true,
                legend: {
                    position: 'left',
                },
                title: {
                    display: true,
                    text: "Current task status",
                    position: "bottom"

                },
                animation: {
                    animateScale: true,
                    animateRotate: true
                },
            }
        });
    }

    private setupAssigneesPie() {

        let elem = document.getElementById("assignees-pie") as any;
        let ctx = elem.getContext("2d");

        let colors = this.assignees.map(() => {
            return this.colors.random[Math.floor(Math.random() * this.colors.random.length)]
        });

        this.assigneesPie = new Chart(ctx, {
            type: "doughnut",
            data: {
                labels: this.assignees.map(x => x.assignee),
                datasets: [
                    {
                        label: "Task status",
                        data: this.assignees.map(x => x.task_count),
                        backgroundColor: colors,
                    }
                ],
            },
            options: {
                responsive: true,
                legend: {
                    position: 'left',
                },
                title: {
                    display: true,
                    text: "Task assignment",
                    position: "bottom"

                },
                animation: {
                    animateScale: true,
                    animateRotate: true
                },
            }
        });
    }

    private getProject() {
        this.apiService.getProject(this.projectId).subscribe((data: any) => {
                this.project = data.content.project;

                this.apiService.getMonitoringSnapshots(60, this.projectId)
                    .subscribe((data: any) => {
                        this.snapshots = data.content.snapshots;
                        this.lastSnapshot = this.snapshots ? this.snapshots.sort((a, b) => {
                            return b.time_stamp - a.time_stamp
                        })[0] : null;

                        this.setupTimeline();
                        this.setupStatusPie();

                        if (!this.snapshots) {
                            return
                        }

                        this.apiService.getAssigneeStats(this.projectId)
                            .subscribe((data: any) => {
                                this.assignees = data.content.assignees;
                                this.setupAssigneesPie();
                            });
                    })
            },
            error => {
                this.translate.get("messenger.unauthorized").subscribe(t =>
                    this.messenger.show(t))
            })
    }

    resetFailedTasks() {
        this.dialog.open(AreYouSureComponent, {
            width: '250px',
        }).afterClosed().subscribe(result => {
            if (result) {
                alert("yes")
            }
        });
    }

    pauseProject() {
        this.dialog.open(AreYouSureComponent, {
            width: '250px',
        }).afterClosed().subscribe(result => {
            if (result) {
                this.project.paused = true;
                this.apiService.updateProject(this.project).subscribe(() => {
                    this.translate.get("messenger.acknowledged").subscribe(t =>
                        this.messenger.show(t))
                })
            }
        });
    }

    hardReset() {

    }
}
