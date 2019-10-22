import { HttpClientModule } from '@angular/common/http';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { GoogleChartsModule } from 'angular-google-charts';
import { CookieService } from 'ngx-cookie-service';
import { AppComponent } from './app.component';
import { ROUTING } from "./app.routing";
import { CapacityGraphModule } from './modules/capacity-graph/capacity-graph.module';
import { ChangepasswordModule } from './modules/changepassword/changepassword.module';
import { LogicalGroupModule } from './modules/logical-group/logical-group.module';
import { LoginModule } from './modules/login/login.module';
import { LogoutModule } from './modules/logout/logout.module';
import { OptionsModule } from './modules/options/options.module';
import { TopoGraphModule } from './modules/topo-graph/modules';
import { TopologyGraphModule } from './modules/topologyGraph/modules';

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
    CapacityGraphModule,
    TopologyGraphModule,
    TopoGraphModule,
    LoginModule,
    LogoutModule,
    LogicalGroupModule,
    ChangepasswordModule,
    OptionsModule,
    GoogleChartsModule.forRoot()
  ],
  providers: [CookieService],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  bootstrap: [AppComponent]
})
export class AppModule { }
