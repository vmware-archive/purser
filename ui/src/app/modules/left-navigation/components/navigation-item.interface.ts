export interface NavigationItemsInterface {
    name: string,
    displayText: any,
    routerLink: string
    visible: boolean
    iconShape?: string,
    activeFor?: Array<string>,
    featureFlag?: string,
    isBeta?: boolean,
    queryParams?: any,
    childItems?: Array<NavigationItemsInterface>
}