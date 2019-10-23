import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { AppComponent } from '../../../app.component';
import { LoginService } from '../services/login.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  public form: any = {};
  public LOGIN_STATUS = "wait";
  ngOnInit() {
    this.LOGIN_STATUS = "wait";
    this.appComponent.IS_LOGEDIN = false;
  }

  constructor(private router: Router, private loginService: LoginService, private appComponent: AppComponent) { }

  public submitLogin() {
    var credentials = JSON.stringify(this.form);
    let observableEntity: Observable<any> = this.loginService.sendLoginCredential(credentials);
    observableEntity.subscribe((response) => {
      this.LOGIN_STATUS = "success";
      this.appComponent.IS_LOGEDIN = true;
      this.router.navigateByUrl('/group');
    }, (err) => {
      this.LOGIN_STATUS = "wrong";
      this.appComponent.IS_LOGEDIN = false;
    });
  }
}