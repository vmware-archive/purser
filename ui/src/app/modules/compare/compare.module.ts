import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClarityModule, ClrFormsNextModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { CompareCloudsComponent } from './components/compare-clouds/compare-clouds.component';

@NgModule({
  imports: [
    CommonModule,
    ClarityModule,
    ClrFormsNextModule,
    FormsModule,
    RouterModule
  ],
  declarations: [CompareCloudsComponent],
})
export class CompareModule { }