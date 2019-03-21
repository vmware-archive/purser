import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BACKEND_URL } from '../../../app.component'

@Injectable()
export class LogicalGroupService {
    constructor(private http: HttpClient) {

    }

    public getLogicalGroupData(name?) {
        let _devUrl: string = './json/logicalGroup.json';
        let _url: string = BACKEND_URL + 'groups';

        if (name) {
            _url = _url + '?name=' + name;
            _devUrl = './json/logicalGroup1.json'; //testing purpose
        }

        //console.log(_url);

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json',
            withCredentials: true,
        });
    }

    public deleteCustomGroup(name) {
        let _url: string = BACKEND_URL + 'group/delete?name=' + name;
        return this.http.post(_url, null, {withCredentials: true})
    }

    public createCustomGroup(groupDef) {
        let _url: string = BACKEND_URL + 'group/create';
        return this.http.post(_url, groupDef, {withCredentials: true})
    }
}