import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { HttpClient } from  "@angular/common/http";
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { ChangepasswordService } from '../services/changepassword.service';

@Component({
    selector: 'app-changepassword',
    templateUrl: './changepassword.component.html',
    styleUrls: ['./changepassword.component.scss']
})
export class ChangepasswordComponent implements OnInit {
    public form: any = {};
    public notMatch = false;
    public LOGIN_STATUS = "wait";
    ngOnInit() {
        this.notMatch = false;
        this.LOGIN_STATUS = "wait";
    }

    constructor(private router: Router, private changepasswordService: ChangepasswordService) { }

    public submitChangePassword() {
        if (this.form.newPassword != this.form.newPasswordCheck) {
            this.notMatch = true;
            return;
        }
        this.notMatch = false;
        var credentials = JSON.stringify(this.form);
        let observableEntity: Observable<any> = this.changepasswordService.sendLoginCredentials(credentials);
        observableEntity.subscribe((response) => {
            this.LOGIN_STATUS = "success";
        }, (err) => {
            this.LOGIN_STATUS = "wrong";
        });
    }
}