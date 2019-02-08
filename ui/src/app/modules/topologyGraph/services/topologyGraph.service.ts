import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class TopologyGraphService {
    constructor(private http: HttpClient) {

    }

    public getNodes(serviceName) {
        let _devUrl: string = './json/nodes.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'nodes';
        if (serviceName && serviceName !== 'ALL') {
            _url = _url + '?service=' + serviceName;
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }

    public getEdges(serviceName) {
        let _devUrl: string = './json/edges.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'edges';
        if (serviceName && serviceName !== 'ALL') {
            _url = _url + '?service=' + serviceName;
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }

    public getServiceList() {
        let _devUrl: string = './json/serviceList.json';
        let base_url: string = window.location.protocol + '//' + window.location.host.split(':')[0] + ':30300/';
        let _url: string = base_url + 'services';

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }
}