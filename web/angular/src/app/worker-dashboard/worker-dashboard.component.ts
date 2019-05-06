import {Component, OnInit} from '@angular/core';
import {ApiService} from '../api.service';

import {Chart} from 'chart.js';
import {AuthService} from "../auth.service";
import {Worker} from "../models/worker";
import {TranslateService} from "@ngx-translate/core";
import {MessengerService} from "../messenger.service";

@Component({
    selector: 'app-worker-dashboard',
    templateUrl: './worker-dashboard.component.html',
    styleUrls: ['./worker-dashboard.component.css']
})
export class WorkerDashboardComponent implements OnInit {

    private chart: Chart;
    workers: Worker[];
    workerInfo: Worker;

    constructor(private apiService: ApiService,
                private translate: TranslateService,
                private messenger: MessengerService,
                public authService: AuthService) {
    }

    ngOnInit() {
        this.setupChart();
        this.refresh();
    }

    public togglePaused(w: Worker) {

        this.workerInfo = undefined;

        this.apiService.workerSetPaused(w.id, !w.paused)
            .subscribe(() => {
                this.refresh();
                this.translate.get('perms.set').subscribe(t => this.messenger.show(t));
            });
    }

    public getInfo(w: Worker) {

        if (this.workerInfo && this.workerInfo.id == w.id) {
            this.workerInfo = undefined;
            return
        }

        this.apiService.getWorker(w.id)
            .subscribe(data => {
                this.workerInfo = data['content']['worker'];
            });
    }

    public refresh() {
        this.apiService.getWorkerStats()
            .subscribe(data => {
                this.updateChart(data['content']['stats']);
                this.workers = data['content']['stats'].sort((a, b) =>
                    (a.alias > b.alias) ? 1 : -1);
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

        data = data
            .filter(w => !w.alias.startsWith("$"))
            .sort((a, b) => b.closed_task_count - a.closed_task_count);

        this.chart.data.labels = data.map(w => w.alias);
        this.chart.data.datasets = [{
            data: data.map(w => w.closed_task_count),
            backgroundColor: '#FF3D00'
        }];
        this.chart.update();
    }
}
