import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { HomeComponent } from './components/home.component';
import { TopoGraphModule } from '../topo-graph/modules';
import { TopologyGraphModule } from '../topologyGraph/modules';
import { CapacityGraphModule } from '../capacity-graph/capacity-graph.module';
import { LogicalGroupModule } from '../logical-group/logical-group.module'

@NgModule({
    imports: [
        ClarityModule,
        FormsModule,
        CommonModule,
        TopoGraphModule,
        TopologyGraphModule,
        CapacityGraphModule,
        LogicalGroupModule
    ],
    declarations: [HomeComponent]
})
export class HomeModule { }
