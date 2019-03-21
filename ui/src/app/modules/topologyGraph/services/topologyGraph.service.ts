import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BACKEND_URL } from '../../../app.component'

@Injectable()
export class TopologyGraphService {
    constructor(private http: HttpClient) {

    }

    public getNodes(serviceName) {
        let _devUrl: string = './json/nodes.json';
        let _url: string = BACKEND_URL + 'nodes';
        if (serviceName && serviceName !== 'ALL') {
            _url = _url + '?service=' + serviceName;
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json',
            withCredentials: true
        });
    }

    public getEdges(serviceName) {
        let _devUrl: string = './json/edges.json';
        let _url: string = BACKEND_URL + 'edges';
        if (serviceName && serviceName !== 'ALL') {
            _url = _url + '?service=' + serviceName;
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json',
            withCredentials: true
        });
    }

    public getServiceList() {
        let _devUrl: string = './json/serviceList.json';
        let _url: string = BACKEND_URL + 'services';

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json',
            withCredentials: true
        });
    }
}