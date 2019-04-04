import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BACKEND_URL } from '../../../app.component';
import { AppComponent } from '../../../app.component';

@Component({
    selector: 'app-options',
    templateUrl: './options.component.html',
    styleUrls: ['./options.component.scss']
})
export class OptionsComponent implements OnInit {
    public SYNC_STATUS = "wait";
    ngOnInit() {
        this.SYNC_STATUS = "wait";
    }

    constructor(private router: Router, private http: HttpClient, private appComponent: AppComponent) { }

    public initiateSync() {
        let syncURL = BACKEND_URL + 'sync';
        const syncOptions = {
            withCredentials: true
        };
        this.http.get(syncURL, syncOptions).subscribe((response) => {
            this.SYNC_STATUS = "requested";
            console.log("sync status", this.SYNC_STATUS);
            }, (err)  => { 
                console.log("Error", err);
                this.SYNC_STATUS = "failed";
                console.log("sync request status", this.SYNC_STATUS);
            }
        );
    }

}