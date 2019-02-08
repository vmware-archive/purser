import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class LogicalGroupService {
    constructor(private http: HttpClient) {

    }

    public getLogicalGroupData(name?) {
        let _devUrl: string = './json/logicalGroup.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'groups';


        if (name) {
            _url = _url + '?name=' + name;
            _devUrl = './json/logicalGroup1.json'; //testing purpose
        }

        //console.log(_url);

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }
}