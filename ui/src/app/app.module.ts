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
import { HomeModule } from './modules/home/home.module';

import { GoogleChartsModule } from 'angular-google-charts';

@NgModule({
    declarations: [
        AppComponent,
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
        LeftNavigationModule,
        HomeModule,
        GoogleChartsModule.forRoot()
    ],
    providers: [],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    bootstrap: [AppComponent]
})
export class AppModule { }
