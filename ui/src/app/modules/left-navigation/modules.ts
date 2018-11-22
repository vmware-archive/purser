import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClarityModule } from '@clr/angular';
import { Routes, RouterModule } from '@angular/router';
import { LeftNavigationComponent } from './components/left-navigation.component';

@NgModule({
    imports: [RouterModule, CommonModule, ClarityModule],
    declarations: [LeftNavigationComponent],
    exports: [LeftNavigationComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class LeftNavigationModule {

}