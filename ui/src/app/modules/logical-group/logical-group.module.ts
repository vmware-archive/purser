import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { GoogleChartsModule } from 'angular-google-charts';
import { LogicalGroupComponent } from './components/logical-group.component';
import { LogicalGroupService } from './services/logical-group.service';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule, GoogleChartsModule
    ],
    exports: [LogicalGroupComponent],
    declarations: [LogicalGroupComponent],
    providers: [LogicalGroupService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class LogicalGroupModule { }
