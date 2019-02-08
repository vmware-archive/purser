import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class CapacityGraphService {
    constructor(private http: HttpClient) {

    }

    public getCapacityData(view?, type?, name?) {
        let _devUrl: string = './json/capacity.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'metrics';

        if (type) {
            _url = _url + '/' + type;
        }

        if (view && !name) {
            _url = _url + '?view=physical';
        }

        if (name) {
            _url = _url + '?name=' + name;
            _devUrl = './json/capacity1.json'; //testing purpose
        }

        //console.log(_url);

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }
}