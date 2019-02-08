/*Framework imports, 3rd party imports */
import { ModuleWithProviders } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './modules/home/components/home.component';
import { LogicalGroupComponent } from './modules/logical-group/components/logical-group.component'
import { TopologyGraphComponent } from './modules/topologyGraph/components/topologyGraph.component'
import { TopoGraphComponent } from './modules/topo-graph/components/topo-graph.component'
import { CapactiyGraphComponent } from './modules/capacity-graph/components/capactiy-graph.component'

export const ROUTES: Routes = [
    { path: 'group', component: LogicalGroupComponent },
    { path: 'network', component: TopologyGraphComponent },
    { path: 'hierarchy', component: TopoGraphComponent },
    { path: 'capacity', component: CapactiyGraphComponent },
    { path: '**', redirectTo: 'hierarchy', pathMatch: 'full' }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);