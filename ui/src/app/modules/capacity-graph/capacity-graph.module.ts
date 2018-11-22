import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { GoogleChartsModule } from 'angular-google-charts';
import { CapactiyGraphComponent } from './components/capactiy-graph.component';
import { CapacityGraphService } from './services/capacity-graph.service';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule, GoogleChartsModule
    ],
    exports: [CapactiyGraphComponent],
    declarations: [CapactiyGraphComponent],
    providers: [CapacityGraphService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class CapacityGraphModule { }
