import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class LogicalGroupService {
    constructor(private http: HttpClient) {

    }

    public getLogicalGroupData(name?) {
        let _devUrl: string = './json/logicalGroup.json';
        let _url: string = 'http://localhost:3030/groups';


        if (name) {
            _url = _url + '?name=' + name;
            _devUrl = './json/logicalGroup1.json'; //testing purpose
        }

        //console.log(_url);

        return this.http.get(_devUrl, {
            observe: 'body',
            responseType: 'json'
        });
    }
}