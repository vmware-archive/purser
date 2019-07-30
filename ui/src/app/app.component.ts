import { Component, OnInit} from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router, RouterEvent, NavigationStart, NavigationEnd, NavigationCancel, NavigationError } from '@angular/router';
import { MCommon } from './common/messages/common.messages';

//production environment
//const BACKEND_BASE_URL = window.location.protocol + '//' + window.location.host

// development environment
 const BACKEND_BASE_URL = 'http://localhost:3030'

export const BACKEND_URL = BACKEND_BASE_URL + '/api/'
export const BACKEND_AUTH_URL = BACKEND_BASE_URL + '/auth/'

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

    public routeLoading: boolean = false;
    public messages: any = {};
    public IS_LOGEDIN = true;

    constructor(public router: Router) {
        this.messages = {
            'common': MCommon
        }
    }

    private loadApp() {
        this.router.events.subscribe((event: RouterEvent) => {
            this.navigationEventHandler(event);
        });
    }

    private navigationEventHandler(event: RouterEvent): void {
        if (event instanceof NavigationStart) {
            this.routeLoading = true;
        }
        if (event instanceof NavigationEnd) {
            this.routeLoading = false;
        }

        // Set loading state to false in both of the below events to hide the spinner in case a request fails.
        if (event instanceof NavigationCancel) {
            this.routeLoading = false;
        }
        if (event instanceof NavigationError) {
            this.routeLoading = false;
        }
    }

    ngOnInit() {
        this.loadApp();
    }

}
