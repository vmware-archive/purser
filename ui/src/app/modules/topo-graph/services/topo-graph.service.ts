import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class TopoGraphService {
    constructor(private http: HttpClient) {

    }

    public getTopoData(view?, type?, name?) {
        let _devUrl: string = './json/topology.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'hierarchy';

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

        //console.log(_url);

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }
}