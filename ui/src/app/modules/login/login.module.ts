import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { LoginComponent } from './components/login.component';
import { LoginService } from './services/login.service';


@NgModule({
    imports: [
        CommonModule, ClarityModule, FormsModule
    ],
    exports: [LoginComponent],
    declarations: [LoginComponent],
    providers: [LoginService],
    schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class LoginModule { }