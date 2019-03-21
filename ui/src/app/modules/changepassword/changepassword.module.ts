import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { ChangepasswordComponent } from './components/changepassword.component';
import { ChangepasswordService } from './services/changepassword.service';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule
    ],
    exports: [ChangepasswordComponent],
    declarations: [ChangepasswordComponent],
    providers: [ChangepasswordService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class ChangepasswordModule { }