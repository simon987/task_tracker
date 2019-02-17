import {Component, OnInit, ViewChild} from '@angular/core';
import {ApiService} from "../api.service";
import {getLogLevel, LogEntry} from "../models/logentry";

import _ from "lodash"
import * as moment from "moment";
import {MatButtonToggleChange, MatPaginator, MatSort, MatTableDataSource} from "@angular/material";

@Component({
    selector: 'app-logs',
    templateUrl: './logs.component.html',
    styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit {

    logs: LogEntry[] = [];
    data: MatTableDataSource<LogEntry>;
    filterLevel: number = 1;
    logsCols: string[] = ["level", "timestamp", "message", "data"];

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    constructor(private apiService: ApiService) {
        this.data = new MatTableDataSource<LogEntry>(this.logs)
    }

    ngOnInit() {
        this.data.paginator = this.paginator;
        this.data.sort = this.sort;
    }

    applyFilter(filter: string) {
        this.data.filter = filter.trim().toLowerCase();
    }

    filterLevelChange(event: MatButtonToggleChange) {
        this.filterLevel = Number(event.value);
        this.getLogs(Number(event.value))
    }

    public refresh() {
        this.getLogs(this.filterLevel)
    }

    private getLogs(level: number) {
        this.apiService.getLogs(level).subscribe(
            data => {
                this.data.data = _.map(data["content"]["logs"], (entry) => {
                    return <LogEntry>{
                        message: entry.message,
                        timestamp: moment.unix(entry.timestamp).format("YYYY-MM-DD HH:mm:ss"),
                        data: JSON.stringify(JSON.parse(entry.data), null, 2),
                        level: getLogLevel(entry.level),
                    }
                });
            }
        )
    }
}

