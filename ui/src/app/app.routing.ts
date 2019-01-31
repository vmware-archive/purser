/*Framework imports, 3rd party imports */
import { ModuleWithProviders } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './modules/home/components/home.component';
import { LogicalGroupComponent } from './modules/logical-group/components/logical-group.component'

export const ROUTES: Routes = [
    { path: 'home', component: HomeComponent },
    { path: 'group', component: LogicalGroupComponent },
    { path: '**', redirectTo: 'home', pathMatch: 'full' }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);