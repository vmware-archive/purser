import { Component } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { MCommon } from '../../../common/messages/common.messages';
import { MLeftNav } from '../../../common/messages/left-navigation.messages';
import { NavigationItemsInterface } from './navigation-item.interface';
import { NavigationItems } from './navigation-items';


@Component({
    selector: 'left-navigation',
    // encapsulation: ViewEncapsulation.None //Default is Emulated
    styleUrls: ['./left-navigation.component.scss'],
    templateUrl: './left-navigation.component.html'
})
export class LeftNavigationComponent {

    public messages: any;
    private isClosed: boolean = false;
    public isCollapsible: boolean = true;
    public collapsed: boolean = true;
    public oldCollapsed: boolean = true;
    private currentHash: string = '';
    public navItems: Array<NavigationItemsInterface> = [];
    public isFreeMode: boolean = false;

    constructor(private router: Router) {
        this.messages = {
            'current': MLeftNav,
            'common': MCommon
        }
        let currentHash: string = '';

        router.events.forEach((event) => {
            if (event instanceof NavigationEnd) {
                currentHash = event.url;
                this.changeActiveRoute(currentHash);
            }
        });
        this.navItems = NavigationItems;
    }

    public changeActiveRoute(hash: string): void {
        this.currentHash = '';
        let hIdx;
        let queryParam = '';
        if ((hIdx = hash.indexOf("?")) != -1 || (hIdx = hash.indexOf(";")) != -1) {
            let eqIndx = hash.indexOf('=');
            queryParam = hash.substr(eqIndx + 1);
            hash = hash.substr(0, hIdx);
        }

        this.currentHash = hash;
    }

    public isActive(item) {
        if (item.activeFor && item.activeFor.indexOf(this.currentHash) !== -1) {
            return 'active activeMenu';
        }
        return '';
    }

    public isRootActive(item) {
        if (this.collapsed && item.activeFor && item.activeFor.indexOf(this.currentHash) !== -1) {
            return 'active activeMenu';
        }
        return '';
    }

    public navigateTo(item) {
        if (item.routerLink) {
            if (item.queryParams) {
                this.router.navigate([item.routerLink], { queryParams: item.queryParams });
            } else {
                this.router.navigate([item.routerLink]);
            }
        }
    }

    public checkVisibility(item) {
        if (item.visible) {
            return true
        }
    }

    public isLocked(item) {
        return false;
    }

    ngDoCheck() {
        if (this.oldCollapsed !== this.collapsed) {
            this.oldCollapsed = this.collapsed;
            setTimeout(() => {
                window.dispatchEvent(new Event('resize')); //manually dispatch the window resize event after few millisecond to resize the chart
            }, 100)

        }
    }
}
