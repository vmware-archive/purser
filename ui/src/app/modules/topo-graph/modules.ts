import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { Routes, RouterModule } from '@angular/router';
import { TopoGraphComponent } from './components/topo-graph.component';
import { TopoGraphService } from './services/topo-graph.service';
import { GoogleChartsModule } from 'angular-google-charts';

@NgModule({
    imports: [RouterModule, CommonModule, ClarityModule, FormsModule, GoogleChartsModule],
    declarations: [TopoGraphComponent],
    exports: [TopoGraphComponent],
    providers: [TopoGraphService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class TopoGraphModule {

}