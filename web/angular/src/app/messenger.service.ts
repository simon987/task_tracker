import {Injectable} from '@angular/core';
import {Subject} from "rxjs";
import {MessengerState} from "./messenger/messenger";

@Injectable()
export class MessengerService {

    public messengerSubject = new Subject<MessengerState>();

    show(message: string) {
        this.messengerSubject.next({
            message: message,
            hidden: false,
        })
    }

    hide() {
        this.messengerSubject.next({
            hidden: true,
        })
    }
}
