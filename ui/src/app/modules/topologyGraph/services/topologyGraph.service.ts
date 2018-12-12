import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class TopologyGraphService {
    constructor(private http: HttpClient) {

    }

    public getNodes(serviceName, nodeType) {
        console.log(nodeType)
        let _devUrl: string = './json/nodes.json';
        let _url: string = 'http://localhost:3030/nodes/' + nodeType;

        if (serviceName && serviceName !== 'ALL') {
            _url = _url + '?service=' + serviceName;
        }

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }

    public getEdges(serviceName, nodeType) {
        console.log(nodeType)
        let _devUrl: string = './json/edges.json';
        let _url: string = 'http://localhost:3030/edges/' + nodeType;

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
        let _url: string = 'http://localhost:3030/services';

        return this.http.get(_url, {
            observe: 'body',
            responseType: 'json'
        });
    }
}