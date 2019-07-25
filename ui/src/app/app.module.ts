import { BrowserModule } from '@angular/platform-browser';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';


import { AppComponent } from './app.component';
import { ROUTING } from "./app.routing";

import { TopologyGraphModule } from './modules/topologyGraph/modules';
import { TopoGraphModule } from './modules/topo-graph/modules';
import { LeftNavigationModule } from './modules/left-navigation/modules';
import { LoginModule } from './modules/login/login.module';
import { LogoutModule } from './modules/logout/logout.module';
import { OptionsModule } from './modules/options/options.module';
import { ChangepasswordModule } from './modules/changepassword/changepassword.module';
import { LogicalGroupModule } from './modules/logical-group/logical-group.module'
import { HomeModule } from './modules/home/home.module';
import { CookieService } from 'ngx-cookie-service';

import { GoogleChartsModule } from 'angular-google-charts';
import { WhatIfComponent } from './what-if/what-if.component';
import { MigrateComponent } from './migrate/migrate.component';
import { FormsModule } from '@angular/forms';

@NgModule({
    declarations: [
        AppComponent,
        WhatIfComponent,
        MigrateComponent,
    ],
    imports: [
        BrowserModule,
        ClarityModule,
        BrowserAnimationsModule,
        RouterModule,
        HttpClientModule,
        ROUTING,
        TopologyGraphModule,
        TopoGraphModule,
        LoginModule,
        LogoutModule,
        LogicalGroupModule,
        ChangepasswordModule,
        LeftNavigationModule,
        HomeModule,
        OptionsModule,
        FormsModule,
        GoogleChartsModule.forRoot()
    ],
    providers: [ CookieService ],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    bootstrap: [AppComponent]
})
export class AppModule { }
