import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {Project} from "../models/project";
import {ActivatedRoute} from "@angular/router";

import {Chart, ChartData, Point} from "chart.js";


@Component({
    selector: 'app-project-dashboard',
    templateUrl: './project-dashboard.component.html',
    styleUrls: ['./project-dashboard.component.css']
})
export class ProjectDashboardComponent implements OnInit {

    private projectId;
    project: Project;
    private timeline: Chart;
    private statusPie: Chart;
    private assigneesPir: Chart;

    private colors = {
        new: "#76FF03",
        failed: "#FF3D00",
        closed: "#E0E0E0",
        awaiting: "#FFB74D"
    };

    tmpLabels = [];
    tmpNew = [];
    tmpFailed = [];
    tmpClosed = [];
    tmpAwaiting = [];

    constructor(private apiService: ApiService, private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getProject();
        });


        let n = 40;
        for (let i = 0; i < n; i++) {
            this.tmpLabels.push((1549501926 + 600 * i) * 1000);
            this.tmpNew.push(Math.ceil(Math.random() * 30))
            this.tmpClosed.push(Math.ceil(Math.random() * 100))
            this.tmpFailed.push(Math.ceil(Math.random() * 13))
            this.tmpAwaiting.push(Math.ceil(Math.random() * 40))
        }

        this.setupTimeline();
        this.setupStatusPie();
        this.setupAssigneesPie();
    }

    private setupTimeline() {
        let elem = document.getElementById("timeline") as any;
        let ctx = elem.getContext("2d");

        this.timeline = new Chart(ctx, {
            type: "bar",
            data: {
                labels: this.tmpLabels,
                datasets: [
                    {
                        label: "New",
                        type: "line",
                        fill: false,
                        borderColor: this.colors.new,
                        backgroundColor: this.colors.new,
                        data: this.tmpNew,
                        pointRadius: 0,
                        lineTension: 0.2,
                    },
                    {
                        label: "Failed",
                        type: "line",
                        fill: false,
                        borderColor: this.colors.failed,
                        backgroundColor: this.colors.failed,
                        data: this.tmpFailed,
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
                        data: this.tmpClosed,
                        lineTension: 0.2,
                    },
                    {
                        label: "Awaiting verification",
                        type: "line",
                        fill: false,
                        borderColor: this.colors.awaiting,
                        backgroundColor: this.colors.awaiting,
                        data: this.tmpAwaiting,
                        pointRadius: 0,
                        lineTension: 0.2,
                    },
                ],
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
                        time: {
                            unit: "minute",
                            unitStepSize: 10,
                        }
                    }]
                },
                tooltips: {
                    enabled: true,
                    intersect: false,
                    mode: "index",
                    position: "nearest",
                },
            }
        })
    }

    private setupStatusPie() {

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
                            10,
                            24,
                            301,
                            90,
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
                }
            }
        });
    }

    private setupAssigneesPie() {

        let elem = document.getElementById("assignees-pie") as any;
        let ctx = elem.getContext("2d");

        this.statusPie = new Chart(ctx, {
            type: "doughnut",
            data: {
                labels: [
                    "marc",
                    "simon",
                    "bernie",
                    "natasha",
                ],
                datasets: [
                    {
                        label: "Task status",
                        data: [
                            10,
                            24,
                            1,
                            23,
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
                    text: "Task assignment",
                    position: "bottom"

                },
                animation: {
                    animateScale: true,
                    animateRotate: true
                }
            }
        });
    }

    private getProject() {
        this.apiService.getProject(this.projectId).subscribe(data => {
            this.project = <Project>{
                id: data["project"]["id"],
                name: data["project"]["name"],
                clone_url: data["project"]["clone_url"],
                git_repo: data["project"]["git_repo"],
                motd: data["project"]["motd"],
                priority: data["project"]["priority"],
                version: data["project"]["version"],
                public: data["project"]["public"],
            }
        })
    }
}
