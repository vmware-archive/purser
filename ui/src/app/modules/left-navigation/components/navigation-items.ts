import { MLeftNav } from '../../../common/messages/left-navigation.messages';
import { NavigationItemsInterface } from './navigation-item.interface';

export const NavigationItems: Array<NavigationItemsInterface> = [
    {
        name: 'Home',
        displayText: MLeftNav.homeText,
        routerLink: '/home',
        visible: true,
        iconShape: 'dashboard',
        activeFor: ['/home']
    }
]