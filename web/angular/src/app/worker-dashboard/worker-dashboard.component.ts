import {Component, OnInit} from '@angular/core';
import {ApiService} from '../api.service';

import {Chart} from 'chart.js';

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
        this.refresh();
    }

    public refresh() {
        this.apiService.getWorkerStats()
            .subscribe(data => {
                this.updateChart(data['content']['stats']);
                }
            );
    }

    private setupChart() {

        const elem = document.getElementById('worker-stats') as any;
        const ctx = elem.getContext('2d');

        this.chart = new Chart(ctx, {
            type: 'bar',
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
        });
    }

    private updateChart(data) {

        data = data.sort((a, b) => b.closed_task_count - a.closed_task_count);

        this.chart.data.labels = data.map(w => w.alias);
        this.chart.data.datasets = [{
            data: data.map(w => w.closed_task_count),
            backgroundColor: '#FF3D00'
        }];
        this.chart.update();
    }
}
