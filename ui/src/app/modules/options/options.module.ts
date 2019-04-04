import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { OptionsComponent } from './components/options.component';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule
    ],
    exports: [OptionsComponent],
    declarations: [OptionsComponent],
    providers: [],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class OptionsModule { }