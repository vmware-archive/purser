import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PlanInfraComponent } from './components/plan-infra/plan-infra.component';
import { ClarityModule, ClrFormsNextModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

@NgModule({
  imports: [
    CommonModule,
    ClarityModule,
    ClrFormsNextModule,
    FormsModule,
    RouterModule
  ],
  declarations: [PlanInfraComponent]
})
export class PlanInfraModule { }
