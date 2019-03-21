import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BACKEND_URL } from '../../../app.component';
import { CookieService } from 'ngx-cookie-service';

@Injectable()
export class TopoGraphService {
    constructor(private http: HttpClient, private cookieService: CookieService) {

    }

    public getTopoData(view?, type?, name?) {
        let _devUrl: string = './json/topology.json';
        let _url: string = BACKEND_URL + 'hierarchy';

        if (type) {
            _url = _url + '/' + type;
        }

        if (view && !name) {
            _url = _url + '?view=physical';
        }

        if (name) {
            _url = _url + '?name=' + name;
            _devUrl = './json/topology1.json'; //testing purpose
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json',
            withCredentials: true,
        });
    }
}