import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { BACKEND_URL } from '../../../app.constants';

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

  constructor(private http: HttpClient) { }

  public initiateSync() {
    let syncURL = BACKEND_URL + 'sync';
    const syncOptions = {
      withCredentials: true
    };
    this.http.get(syncURL, syncOptions).subscribe((response) => {
      this.SYNC_STATUS = "requested";
      console.log("sync status", this.SYNC_STATUS);
    }, (err) => {
      console.log("Error", err);
      this.SYNC_STATUS = "failed";
      console.log("sync request status", this.SYNC_STATUS);
    });
  }

}