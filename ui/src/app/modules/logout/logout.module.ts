import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { LogoutComponent } from './components/logout.component';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule
    ],
    exports: [LogoutComponent],
    declarations: [LogoutComponent],
    providers: [],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class LogoutModule { }