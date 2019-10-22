import { Component, OnInit } from '@angular/core';
import { NavigationCancel, NavigationEnd, NavigationError, NavigationStart, Router, RouterEvent } from '@angular/router';
import { MCommon } from './common/messages/common.messages';

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
