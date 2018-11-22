/*Framework imports, 3rd party imports */
import { ModuleWithProviders } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TopologyGraphComponent } from './modules/topologyGraph/components/topologyGraph.component';
import { TopoGraphComponent } from './modules/topo-graph/components/topo-graph.component';
import { HomeComponent } from './modules/home/components/home.component';

export const ROUTES: Routes = [
    { path: 'home', component: HomeComponent },
    { path: '**', redirectTo: 'home', pathMatch: 'full' }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);