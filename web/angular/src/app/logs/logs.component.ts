import {Component, OnInit, ViewChild} from '@angular/core';
import {ApiService} from "../api.service";

import _ from "lodash"
import * as moment from "moment";
import {MatPaginator, MatSort, MatTableDataSource} from "@angular/material";

@Component({
    selector: 'app-logs',
    templateUrl: './logs.component.html',
    styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit {

    logs: LogEntry[] = [];
    data: MatTableDataSource<LogEntry>;
    logsCols: string[] = ["level", "timestamp", "message", "data"];

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(private apiService: ApiService) {
        this.data = new MatTableDataSource<LogEntry>(this.logs)
    }

    ngOnInit() {
        this.getLogs();

        this.data.paginator = this.paginator;
        this.data.sort = this.sort;
        // interval(5000).subscribe(() => {
        //     this.getLogs();
        // })
    }

    applyFilter(filter: string) {
        this.data.filter = filter.trim().toLowerCase();
    }

    private getLogs() {
        this.apiService.getLogs().subscribe(
            data => {
                this.data.data = _.map(data["logs"], (entry) => {
                    return <LogEntry>{
                        message: entry.message,
                        timestamp: moment.unix(entry.timestamp).toISOString(),
                        data: JSON.stringify(JSON.parse(entry.data), null, 2),
                        level: entry.level
                    }
                });
            }
        )
    }
}

