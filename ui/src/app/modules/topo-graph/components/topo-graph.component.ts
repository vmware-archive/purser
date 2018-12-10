import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { TopoGraphService } from '../services/topo-graph.service';

const STATUS_WAIT = 'WAIT',
    STATUS_READY = 'READY',
    STATUS_NODATA = 'NO_DATA';

@Component({
    selector: 'app-topo-graph',
    templateUrl: './topo-graph.component.html',
    styleUrls: ['./topo-graph.component.scss']
})
export class TopoGraphComponent implements OnInit {
    //PUBLIC
    public TOPO_STATUS = STATUS_WAIT;
    public graphData = [];
    public colNames = ['Child', 'Parent', 'Tooltip']
    public chartOptions = {
        nodeClass: 'customNode',
        allowHtml: true
    };
    public physicalView: boolean = false;
    public legendArr: any = [
        {
            displayText: 'namespace',
            color: 'red'
        },
        {
            displayText: 'service',
            color: 'yellow'
        },
        {
            displayText: 'pod',
            color: 'green'
        },
        {
            displayText: 'container',
            color: 'blue'
        },
        {
            displayText: 'process',
            color: 'orange'
        },
        {
            displayText: 'cluster',
            color: 'orangered'
        },
        {
            displayText: 'deployment',
            color: 'purple'
        },
        {
            displayText: 'replicaset',
            color: 'palevioletred'
        },
        {
            displayText: 'node',
            color: 'royalblue'
        },
        {
            displayText: 'daemonset',
            color: 'brown'
        },
        {
            displayText: 'job',
            color: 'black'
        },
        {
            displayText: 'statefulset',
            color: 'goldenrod'
        }
    ];
    public filterItems: any = [];
    public selectedFilterItem: string = 'select';

    //PRIVATE
    private orgTopoData: any = {};
    private topoData: any = {};
    private keysToConsider: any = ['service', 'pod', 'container', 'process', 'cluster', 'namespace', 'deployment', 'replicaset', 'node', 'daemonset', 'job', 'statefulset', 'children'];
    private uniqNames: any = [];

    constructor(private router: Router, private topoService: TopoGraphService) { }

    private getTopoData() {
        this.graphData = [];
        this.uniqNames = [];
        let observableEntity: Observable<any> = this.topoService.getTopoData(this.physicalView);
        this.TOPO_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
                return;
            }
            this.topoData = response && response.data || {};
            this.orgTopoData = JSON.parse(JSON.stringify(this.topoData));
            this.constructData(this.topoData);
            //console.log(this.topoData);
        }, (err) => {
            this.TOPO_STATUS = STATUS_NODATA;
        });
    }

    private getIcon(type) {
        if (!type) {
            return 'host';
        }
        switch (type) {
            case 'service':
                return 'cluster';
            case 'pod':
                return 'storage';
            case 'container':
                return 'host';
            default:
                return 'host';

        }
    }

    private pushToGraphData(item, parent) {
        let eachRow = [];
        let iconShape = this.getIcon(item.type);
        let parentName = item.name === parent.name ? parent.type : parent.name;
        eachRow.push({ v: item.name, f: '<span class="' + item.type + '">' + item.name + '</span>', t: item.type });
        eachRow.push(parentName);
        eachRow.push(item.type);
        if (this.uniqNames.indexOf(item.name) === -1) {
            this.graphData.push(eachRow);
            this.uniqNames.push(item.name);
        }
    }

    /*private traverse(data){
        for (let key in data) {
            if (this.keysToConsider.indexOf(key) > -1) {
                for (let item of data[key]) {
                    this.pushToGraphData(item, item);
                }
            }
        }
    }*/

    private collectFilterItems(item) {
        this.filterItems.push(item.name);
    }

    private constructData(topoData) {
        let data = JSON.parse(JSON.stringify(topoData));
        for (let key in data) {
            if (this.keysToConsider.indexOf(key) > -1) {
                this.filterItems = [];
                for (let item of data[key]) {
                    this.collectFilterItems(item);
                    this.pushToGraphData(item, data);
                    /*for (let subKey in item) {
                        if (this.keysToConsider.indexOf(subKey) > -1) {
                            for (let subItem of item[subKey]) {
                                this.pushToGraphData(subItem, item);
                                for (let supSubKey in subItem) {
                                    if (this.keysToConsider.indexOf(supSubKey) > -1) {
                                        for (let supSubItem of subItem[supSubKey]) {
                                            this.pushToGraphData(supSubItem, subItem);
                                        }
                                    }
                                }
                            }
                        }
                    }*/
                }
            }
        }
        this.TOPO_STATUS = STATUS_READY;
        //console.log(this.graphData);
    }

    public onSelect(element) {
        if (!element) {
            return;
        }
        if (!element[0]) {
            return;
        }
        let row = element[0].row;
        let selectedItem = this.graphData[row];
        this.getAdditionalData(selectedItem);
    }

    private getAdditionalData(item) {
        this.selectedFilterItem = 'select';
        if (item && item[0] && item[0].v && item[0].t) {
            let name = item[0].v;
            let type = item[0].t;
            let observableEntity: Observable<any> = this.topoService.getTopoData(this.physicalView, type, name);
            this.TOPO_STATUS = STATUS_WAIT;
            observableEntity.subscribe((response) => {
                if (!response) {
                    return;
                }
                let topoData = response && response.data || {};
                this.constructData(topoData);
            }, (err) => {
                this.TOPO_STATUS = STATUS_NODATA;
            });
        } else {
            return;
        }

    }

    public filterItemChange(evt) {
        for (let item of this.graphData) {
            for (let subItem of item) {
                if (subItem && subItem.v && subItem.v === this.selectedFilterItem) {
                    this.uniqNames = [];
                    this.graphData = [];
                    this.getAdditionalData(item);
                }
            }
        }
    }

    public reset() {
        this.TOPO_STATUS = STATUS_WAIT;
        this.graphData = [];
        this.uniqNames = [];
        this.constructData(this.orgTopoData);
    }

    public viewChange() {
        this.getTopoData();
    }

    ngOnInit() {
        this.getTopoData();
    }

}
