import {Component, OnInit} from '@angular/core';
import {AuthService} from "../auth.service";
import {ApiService} from "../api.service";
import {ActivatedRoute} from "@angular/router";
import {TranslateService} from "@ngx-translate/core";
import {MessengerService} from "../messenger.service";

@Component({
    selector: 'app-project-secret',
    templateUrl: './project-secret.component.html',
    styleUrls: ['./project-secret.component.css']
})
export class ProjectSecretComponent implements OnInit {

    secret: string;
    webhookSecret: string;
    projectId: number;

    constructor(private auth: AuthService,
                private apiService: ApiService,
                private translate: TranslateService,
                private messenger: MessengerService,
                private route: ActivatedRoute) {
    }

    ngOnInit() {
        this.route.params.subscribe(params => {
            this.projectId = params["id"];
            this.getSecret();
            this.getWebhookSecret();
        });
    }

    getSecret() {
        this.apiService.getSecret(this.projectId).subscribe(data => {
            this.secret = data["content"]["secret"]
        }, error => {
            this.translate.get("messenger.unauthorized").subscribe(t => this.messenger.show(t))
        })
    }

    getWebhookSecret() {
        this.apiService.getWebhookSecret(this.projectId).subscribe(data => {
            this.webhookSecret = data["content"]["webhook_secret"]
        }, error => {
            this.translate.get("messenger.unauthorized").subscribe(t => this.messenger.show(t))
        })
    }

    onUpdate() {
        this.apiService.setSecret(this.projectId, this.secret).subscribe(data => {
            this.translate.get("secret.ok").subscribe(t => this.messenger.show(t))
        })
    }

    onWebhookUpdate() {
        this.apiService.setWebhookSecret(this.projectId, this.webhookSecret).subscribe(data => {
            this.translate.get("secret.ok").subscribe(t => this.messenger.show(t))
        })
    }

    refresh() {
        this.getWebhookSecret();
        this.getSecret();
    }
}
