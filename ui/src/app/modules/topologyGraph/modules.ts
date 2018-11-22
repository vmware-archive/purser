import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { Routes, RouterModule } from '@angular/router';
import { TopologyGraphComponent } from './components/topologyGraph.component';
import { TopologyGraphService } from './services/topologyGraph.service';

@NgModule({
    imports: [RouterModule, CommonModule, ClarityModule, FormsModule],
    declarations: [TopologyGraphComponent],
    exports: [TopologyGraphComponent],
    providers: [TopologyGraphService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class TopologyGraphModule {

}