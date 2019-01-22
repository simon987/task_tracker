import {Component, OnInit} from '@angular/core';

import * as d3 from "d3"
import * as _ from "lodash"
import {interval} from "rxjs";
import {ApiService} from "../api.service";

@Component({
    selector: 'app-project-dashboard',
    templateUrl: './project-dashboard.component.html',
    styleUrls: ['./project-dashboard.component.css']
})
export class ProjectDashboardComponent implements OnInit {

    projectStats;
    private pieWidth = 360;
    private pieHeight = 360;
    private pieRadius = Math.min(this.pieWidth, this.pieHeight) / 2;
    private pieArc = d3.arc()
        .innerRadius(this.pieRadius / 2)
        .outerRadius(this.pieRadius);

    private pieFun = d3.pie().value((d) => d.count);
    private statusColor = d3.scaleOrdinal().range(['#31a6a2', '#8c2627', '#62f24b']);
    private assigneesColor = d3.scaleOrdinal().range(["", "#AAAAAA"].concat(d3.schemePaired));

    private statusData: any[];
    private assigneesData: any[];

    private newTaskCounts: any[] = [];
    private failedTaskCounts: any[] = [];
    private closedTaskCounts: any[] = [];

    private newTaskPath: any;
    private failedTaskPath: any;
    private closedTaskPath: any;

    private maxY: number = 10;
    private yAxis: any;
    private yScale: any;
    private xScale: any;
    private range: number;
    private line: any;

    private statusPath: any;
    private statusSvg: any;
    private assigneesPath: any;
    private assigneesSvg: any;

    constructor(private apiService: ApiService) {
    }

    setupStatusPieChart() {
        let tooltip = d3.select("#stooltip");

        this.statusSvg = d3.select('#status')
            .append('svg')
            .attr('width', this.pieWidth)
            .attr('height', this.pieHeight)
            .append("g")
            .attr("transform", "translate(" + this.pieRadius + "," + this.pieRadius + ")");

        this.statusPath = this.statusSvg.selectAll("path")
            .data(this.pieFun(this.statusData))
            .enter()
            .append('path')
            .attr('d', this.pieArc)
            .attr('fill', (d) => this.statusColor(d.data.label));

        this.setupToolTip(this.statusPath, tooltip)
    }

    setupAssigneesPieChart() {
        let tooltip = d3.select("#atooltip");

        this.assigneesSvg = d3.select('#assignees')
            .append('svg')
            .attr('width', this.pieWidth)
            .attr('height', this.pieHeight)
            .append("g")
            .attr("transform", "translate(" + this.pieRadius + "," + this.pieRadius + ")");

        this.assigneesPath = this.assigneesSvg.selectAll("path")
            .data(this.pieFun(this.assigneesData))
            .enter()
            .append('path')
            .attr('d', this.pieArc)
            .attr('fill', (d) => this.assigneesColor(d.data.label));

        this.setupToolTip(this.assigneesPath, tooltip)
    }

    setupToolTip(x, tooltip) {
        x.on('mouseover', (d) => {
            let total = d3.sum(this.assigneesData.map((d) => d.count));
            let percent = Math.round(1000 * d.data.count / total) / 10;
            tooltip.select('.label').html(d.data.label);
            tooltip.select('.count').html(d.data.count);
            tooltip.select('.percent').html(percent + '%');
            tooltip.style('display', 'block');
        });
        x.on('mouseout', function () {
            tooltip.style('display', 'none');
        })
    }

    setupLine() {
        let margin = {top: 50, right: 50, bottom: 50, left: 50};
        this.range = 600;
        let width = 750;
        let height = 250;

        this.xScale = d3.scaleLinear()
            .domain([this.range, 0])
            .range([width, 0]);

        this.yScale = d3.scaleLinear()
            .domain([0, this.maxY])
            .range([height, 0]);

        this.line = d3.line()
            .x((d, i) => this.xScale(i))
            .y((d) => this.yScale(d.y))
            .curve(d3.curveMonotoneX);

        let svg = d3.select("#line").append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

        svg.append("defs").append("clipPath")
            .attr("id", "clip")
            .append("rect")
            .attr("width", width)
            .attr("height", height);

        this.newTaskPath = svg
            .append("path")
            .attr("clip-path", "url(#clip)")
            .datum(this.newTaskCounts)
            .attr("class", "line-new")
            .attr("d", this.line);
        this.failedTaskPath = svg
            .append("path")
            .attr("clip-path", "url(#clip)")
            .datum(this.failedTaskCounts)
            .attr("class", "line-failed")
            .attr("d", this.line);
        this.closedTaskPath = svg
            .append("path")
            .attr("clip-path", "url(#clip)")
            .datum(this.closedTaskCounts)
            .attr("class", "line-closed")
            .attr("d", this.line);

        let xAxis = svg.append("g")
            .attr("class", "x-axis")
            .attr("transform", "translate(0," + height + ")");
        xAxis.call(d3.axisBottom(this.xScale).tickFormat((d) => (d - 600) + "s"));

        this.yAxis = svg.append("g")
            .attr("class", "y-axis");
        this.yAxis.call(d3.axisLeft(this.yScale));
    }

    getStats() {
        this.apiService.getProjectStats(2).subscribe((data) => {

            this.projectStats = data["stats"];

            this.updateLine();
            this.updatePie();
        });
    }

    private updateLine() {

        let newVal = {"y": this.projectStats["new_task_count"]};
        let failedVal = {"y": this.projectStats["failed_task_count"]};
        let closedVal = {"y": this.projectStats["closed_task_count"]};

        //Adjust y axis
        this.maxY = Math.max(newVal["y"], this.maxY);
        this.yScale.domain([0, this.maxY]);
        this.yAxis.call(d3.axisLeft(this.yScale));

        this.newTaskPath
            .attr("d", this.line)
            .attr("transform", null);
        this.failedTaskPath
            .attr("d", this.line)
            .attr("transform", null);
        this.closedTaskPath
            .attr("d", this.line)
            .attr("transform", null);

        //remove fist element
        if (this.newTaskCounts.length >= this.range) {
            this.newTaskCounts.shift();
            this.newTaskPath
                .transition()
                .attr("transform", "translate(" + this.xScale(-1) + ")");
        }
        if (this.failedTaskCounts.length >= this.range) {
            this.failedTaskCounts.shift();
            this.failedTaskPath
                .transition()
                .attr("transform", "translate(" + this.xScale(-1) + ")");
        }
        if (this.closedTaskCounts.length >= this.range) {
            this.closedTaskCounts.shift();
            this.closedTaskPath
                .transition()
                .attr("transform", "translate(" + this.xScale(-1) + ")");
        }

        this.newTaskCounts.push(newVal);
        this.failedTaskCounts.push(failedVal);
        this.closedTaskCounts.push(closedVal);
    }

    private updatePie() {
        this.statusData = [
            {label: "New", count: this.projectStats["new_task_count"]},
            {label: "Failed", count: this.projectStats["failed_task_count"]},
            {label: "Closed", count: this.projectStats["closed_task_count"]},
        ];
        this.assigneesData = _.map(this.projectStats["assignees"], (assignedTasks) => {
            return {
                label: assignedTasks["assignee"] == "00000000-0000-0000-0000-000000000000" ? "unassigned" : assignedTasks["assignee"],
                count: assignedTasks["task_count"]
            }
        });

        this.statusSvg.selectAll("path")
            .data(this.pieFun(this.statusData));
        this.statusPath
            .attr('d', this.pieArc)
            .attr('fill', (d) => this.statusColor(d.data.label));
        this.assigneesSvg.selectAll("path")
            .data(this.pieFun(this.assigneesData));
        this.assigneesPath
            .attr('d', this.pieArc)
            .attr('fill', (d) => this.assigneesColor(d.data.label));
    }

    ngOnInit() {
        this.statusData = [
            {label: 'new', count: 0},
            {label: 'failed', count: 0},
            {label: 'closed', count: 0},
        ];
        this.assigneesData = [
            {label: 'null', count: 0},
        ];

        this.setupStatusPieChart();
        this.setupAssigneesPieChart();
        this.setupLine();

        this.getStats();
        interval(1000).subscribe(() => {
            this.getStats()
        })
    }
}
