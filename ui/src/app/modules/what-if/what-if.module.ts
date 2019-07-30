import { NgModule , CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClarityModule, ClrFormsNextModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { WhatIfComponent } from './components/whatif/what-if.component';
import { RouterModule } from '@angular/router';

@NgModule({
  imports: [
    CommonModule,
    ClarityModule,
    ClrFormsNextModule,
    FormsModule,
    RouterModule
  ],
  declarations: [WhatIfComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
  //providers: [AppProfileService],
})
export class WhatIfModule { }
