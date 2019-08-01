/*Framework imports, 3rd party imports */
import { ModuleWithProviders } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LogicalGroupComponent } from './modules/logical-group/components/logical-group.component'
import { TopologyGraphComponent } from './modules/topologyGraph/components/topologyGraph.component'
import { TopoGraphComponent } from './modules/topo-graph/components/topo-graph.component'
import { CapactiyGraphComponent } from './modules/capacity-graph/components/capactiy-graph.component'
import { LoginComponent } from './modules/login/components/login.component'
import { LogoutComponent } from './modules/logout/components/logout.component'
import { OptionsComponent } from './modules/options/components/options.component'
import { ChangepasswordComponent } from './modules/changepassword/components/changepassword.component'
import { WhatIfComponent } from './modules/what-if/components/whatif/what-if.component';
import { CompareCloudsComponent } from './modules/compare/components/compare-clouds/compare-clouds.component';
import { PlanInfraComponent } from './modules/plan-infra/components/plan-infra/plan-infra.component';

export const ROUTES: Routes = [
    { path: 'group', component: LogicalGroupComponent },
    { path: 'network', component: TopologyGraphComponent },
    { path: 'hierarchy', component: TopoGraphComponent },
    { path: 'capacity', component: CapactiyGraphComponent },
    { path: 'login', component: LoginComponent },
    { path: 'logout', component: LogoutComponent },
    { path: 'changepassword', component: ChangepasswordComponent },
    { path: 'options', component: OptionsComponent },
    { path: 'whatif', component: WhatIfComponent},
    { path: 'compare', component : CompareCloudsComponent},
    { path: 'plan/infra', component: PlanInfraComponent},
    { path: '**', redirectTo: 'login', pathMatch: 'full' }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);