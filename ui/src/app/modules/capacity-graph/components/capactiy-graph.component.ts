import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { CapacityGraphService } from '../services/capacity-graph.service';

const STATUS_WAIT = 'WAIT',
    STATUS_READY = 'READY',
    STATUS_NODATA = 'NO_DATA';

@Component({
    selector: 'app-capactiy-graph',
    templateUrl: './capactiy-graph.component.html',
    styleUrls: ['./capactiy-graph.component.scss']
})


export class CapactiyGraphComponent implements OnInit {
    public cpuAllocated = 100.0;
    public cpuCapacity = 100.0;
    public cpuRatio = 100;
    public memoryAllocated = 100.0;
    public memoryCapacity = 100.0;
    public memoryRatio = 100;
    public storageAllocated = 100.0;
    public storageCapacity = 100.0;
    public storageRatio = 100;
    public resourceType = 'Cluster';

    //PUBLIC
    public CAPA_STATUS = STATUS_WAIT;
    public graphData = [];
    public colNames = ['Child', 'Parent', 'Metrics'];
    public chartOptions = {
        nodeClass: 'customNode',
        allowHtml: true,
        animation: {
            startup: true,
            duration: 1000,
            easing: 'out',
        },
        minColor: '#009688',
        midColor: '#f7f7f7',
        maxColor: '#ee8100',
        headerHeight: 40,
    };
    public selectedMetric: string = 'cpu';
    public metricOptions: any = [
        { displayValue: 'CPU', value: 'cpu', units: 'vCPU' },
        { displayValue: 'Memory', value: 'memory', units: 'GB' },
        { displayValue: 'Storage', value: 'storage', units: 'GB' }
        //{ displayValue: 'Network', value: 'network' }
    ];
    public physicalView: boolean = false;
    public rootItem: any = {};
    public filterItems: any = [];
    public selectedFilterItem: string = 'select';

    //PRIVATE
    private orgCapaData: any = {};
    private capaData: any = {};
    private keysToConsider: any = ['service', 'pod', 'container', 'process', 'cluster', 'namespace', 'deployment', 'replicaset', 'node', 'daemonset', 'job', 'statefulset', 'children', 'pv', 'pvc'];
    private uniqNames: any = [];

    constructor(private router: Router, private capacityGraphService: CapacityGraphService) { }

    private getCapacityData() {
        let observableEntity: Observable<any> = this.capacityGraphService.getCapacityData(this.physicalView);
        this.CAPA_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
                return;
            }
            this.capaData = response && response.data || {};
            this.orgCapaData = JSON.parse(JSON.stringify(this.capaData));

            this.constructData(this.capaData);
        }, (err) => {
            this.CAPA_STATUS = STATUS_NODATA;
        });
    }

    private computeAllocationRatios(data) {
        if (!!data.cpuCapacity) {
            this.cpuCapacity = data.cpuCapacity.toFixed(2);
        } else {
            this.cpuCapacity = 0;
        }
        if (!!data.cpuAllocated) {
            this.cpuAllocated = data.cpuAllocated.toFixed(2);
        } else {
            this.cpuAllocated = 0;
        }
        this.cpuRatio = Math.round(this.cpuAllocated * 100 / this.cpuCapacity);

        if (!!data.memoryCapacity) {
            this.memoryCapacity = data.memoryCapacity.toFixed(2);
        } else {
            this.memoryCapacity = 0;
        }
        if (!!data.memoryAllocated) {
            this.memoryAllocated = data.memoryAllocated.toFixed(2);
        } else {
            this.memoryAllocated = 0;
        }
        this.memoryRatio = Math.round(this.memoryAllocated * 100 / this.memoryCapacity);

        if (!!data.storageCapacity) {
            this.storageCapacity = data.storageCapacity.toFixed(2);
        } else {
            this.storageCapacity = 0;
        }
        if (!!data.storageAllocated) {
            this.storageAllocated = data.storageAllocated.toFixed(2);
        } else {
            this.storageAllocated = 0;
        }
        this.storageRatio = Math.round(this.storageAllocated * 100 / this.storageCapacity);

        if (data.type == 'node') {
            this.resourceType = 'Node';
            this.storageAllocated = 0;
            this.storageCapacity = 0;
            this.storageRatio = 0;
        } else {
            if (data.type == 'pv') {
                this.resourceType = 'PersistentVolume';
                this.cpuAllocated = 0;
                this.cpuCapacity = 0;
                this.cpuRatio = 0;
                this.memoryAllocated = 0;
                this.memoryCapacity = 0;
                this.memoryRatio = 0;
            } else {
                this.resourceType = 'Cluster';
            }
        }
    }

    private constructRoot(capaData) {
        this.rootItem = capaData;
        let eachRow = [];
        let rootName = capaData && capaData.name;
        let metricValue = capaData[this.selectedMetric] || 0;
        let metricCostValue = capaData[this.selectedMetric + 'Cost'] || 0;
        if (rootName) {
            eachRow.push({ v: rootName, f: rootName + ', ' + this.selectedMetric + ': ' + metricValue.toFixed(2) + ', ' + this.selectedMetric + ' cost: ' + metricCostValue.toFixed(2) });
            eachRow.push(null);
            eachRow.push(0);
            if (this.uniqNames.indexOf(rootName) === -1) {
                this.graphData.push(eachRow);
                this.uniqNames.push(rootName);
            }
        }
    }

    private pushToGraphData(item, parent) {
        let eachRow = [];
        let parentName = item.name === parent.name ? parent.type : parent.name;
        let metricValue = item[this.selectedMetric] || 0;
        let metricCostValue = item[this.selectedMetric + 'Cost'] || 0;
        eachRow.push({ v: item.name, f: item.name + ', ' + this.selectedMetric + ': ' + metricValue.toFixed(2) + ', ' + this.selectedMetric + ' cost: ' + metricCostValue.toFixed(2), t: item.type });
        eachRow.push(parentName);
        eachRow.push(metricValue);
        if (this.uniqNames.indexOf(item.name) === -1) {
            this.graphData.push(eachRow);
            this.uniqNames.push(item.name);
        }
    }

    private collectFilterItems(item) {
        this.filterItems.push(item.name);
    }

    private constructData(capaData) {
        this.selectedFilterItem = 'select';
        this.graphData = [];
        this.uniqNames = [];
        this.constructRoot(capaData);
        let data = JSON.parse(JSON.stringify(capaData));
        for (let key in data) {
            if (this.keysToConsider.indexOf(key) > -1) {
                this.filterItems = [];
                for (let item of data[key]) {
                    this.collectFilterItems(item);
                    this.pushToGraphData(item, data);
                }
            }
        }

        this.computeAllocationRatios(data);

        this.CAPA_STATUS = STATUS_READY;
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
        let selectedItem = item;
        if (item && item[0] && item[0].v && item[0].t) {
            let name = item[0].v;
            let type = item[0].t;
            let observableEntity: Observable<any> = this.capacityGraphService.getCapacityData(this.physicalView, type, name);
            this.CAPA_STATUS = STATUS_WAIT;
            observableEntity.subscribe((response) => {
                if (!response) {
                    return;
                }
                let capaData = response && response.data || {};
                this.constructData(capaData);
            }, (err) => {
                this.CAPA_STATUS = STATUS_NODATA;
            });
        } else {
            return;
        }
    }

    public filterItemChange(evt) {
        for (let item of this.graphData) {
            for (let subItem of item) {
                if (subItem && subItem.v && subItem.v === this.selectedFilterItem) {
                    this.getAdditionalData(item);
                }
            }
        }
    }

    public metricChange(evt) {
        this.constructData(this.orgCapaData);
    }

    public reset() {
        this.CAPA_STATUS = STATUS_WAIT;
        this.graphData = [];
        this.uniqNames = [];
        this.constructData(this.orgCapaData);
    }

    public viewChange() {
        this.getCapacityData()
    }

    ngOnInit() {
        this.getCapacityData();
    }

}
