import {Component, OnInit} from '@angular/core';
import {ApiService} from "../api.service";

import {Chart} from "chart.js";

@Component({
    selector: 'app-worker-dashboard',
    templateUrl: './worker-dashboard.component.html',
    styleUrls: ['./worker-dashboard.component.css']
})
export class WorkerDashboardComponent implements OnInit {

    private chart: Chart;

    constructor(private apiService: ApiService) {
    }

    ngOnInit() {
        this.setupChart();
        this.refresh()
    }

    private refresh() {
        this.apiService.getWorkerStats()
            .subscribe((data: any) => {
                    this.updateChart(data.stats)
                }
            )
    }

    private setupChart() {

        let elem = document.getElementById("worker-stats") as any;
        let ctx = elem.getContext("2d");

        this.chart = new Chart(ctx, {
            type: "bar",
            data: {
                labels: [],
                datasets: [],
            },
            options: {
                title: {
                    display: false,
                },
                legend: {
                    display: false
                },
                tooltips: {
                    enabled: true,
                },
                responsive: true
            }
        })
    }

    private updateChart(data) {

        this.chart.data.labels = data.map(w => w.alias);
        this.chart.data.datasets = [{
            data: data.map(w => w.closed_task_count),
            backgroundColor: "#FF3D00"
        }];
        this.chart.update();
    }
}
