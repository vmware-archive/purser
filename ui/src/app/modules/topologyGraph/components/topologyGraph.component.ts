import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { DataSet, Network } from 'vis';
import { Observable } from 'rxjs';
import { TopologyGraphService } from '../services/topologyGraph.service';

const STATUS_WAIT = 'WAIT',
    STATUS_READY = 'READY',
    STATUS_NODATA = 'NO_DATA';

@Component({
    selector: 'topology-graph',
    templateUrl: './topologyGraph.component.html',
    styleUrls: ['./topologyGraph.component.scss']
})

export class TopologyGraphComponent implements OnInit {
    private clusterIndex = 0;
    private clusters = [];
    private lastClusterZoomLevel = 0;
    private clusterFactor = 0.9;
    private nodes: any;
    private edges: any;
    private nodesDataset: any;
    private edgesDataset: any;

    public NODE_STATUS = STATUS_WAIT;
    public EDGE_STATUS = STATUS_WAIT;
    public serviceList: any = [];
    public enableClustering: boolean = false;
    public serviceName: string = 'ALL';

    @ViewChild('networkContainer') container: ElementRef;


    data: any = {};
    options = {
        nodes: {
            shape: 'dot',
            size: 16
        },
        physics: {
            enabled: false,
            /*forceAtlas2Based: {
                gravitationalConstant: -26,
                centralGravity: 0.005,
                springLength: 230,
                springConstant: 0.18
            },
            maxVelocity: 146,
            solver: 'forceAtlas2Based',
            timestep: 0.35,
            stabilization: { iterations: 150 }*/
        },
        layout: {
            improvedLayout: false
        }
    };

    network: any;

    constructor(private router: Router, private topologyService: TopologyGraphService) {

    }

    private getServiceList() {
        let observableEntity: Observable<any> = this.topologyService.getServiceList();
        this.NODE_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
                return;
            }
            this.serviceList = response;
        }, (err) => {
        });
    }

    private getNodes(nodeType) {
        if (nodeType == 'pod') {
            this.serviceList = [];
        }
        let observableEntity: Observable<any> = this.topologyService.getNodes(this.serviceName, nodeType);
        this.NODE_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
                this.NODE_STATUS = STATUS_NODATA;
                return;
            }
            this.nodes = response;
            for (let item of this.nodes) {
                if (item.cid && this.serviceList.indexOf(item.cid) === -1) {
                    for (let cid of item.cid) {
                        if (this.serviceList.indexOf(cid) === -1) {
                            this.serviceList.push(cid);
                        }
                    }
                }
            }
            this.NODE_STATUS = STATUS_READY;
            this.initNetwork();
        }, (err) => {
            this.NODE_STATUS = STATUS_NODATA;
        });
    }

    private getEdges(nodeType) {
        let observableEntity: Observable<any> = this.topologyService.getEdges(this.serviceName, nodeType);
        this.EDGE_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
                this.EDGE_STATUS = STATUS_NODATA;
                return;
            }
            this.edges = response;
            this.EDGE_STATUS = STATUS_READY;
            this.initNetwork();
        }, (err) => {
            this.EDGE_STATUS = STATUS_NODATA;
        });
    }

    private initNetwork() {
        let filteredNodes = [];
        let filteredEdges = [];
        // console.log(this.serviceName)
        if (this.EDGE_STATUS === STATUS_READY && this.NODE_STATUS === STATUS_READY) {
            if (this.serviceName && this.serviceName !== 'ALL') {
                let self = this;
                // console.log(this.nodes, "All nodes")
                filteredNodes = this.nodes.filter(function (item) {
                    // console.log(item, "item")
                    return item.label == self.serviceName;
                });
                // console.log(filteredNodes, "After filter")
                let idsArr = [];
                for (let item of filteredNodes) {
                    idsArr.push(item.id);
                }
                // console.log(idsArr, "idsArr");
                // console.log(this.edges, "All edges")
                filteredEdges = this.edges.filter(function (item) {
                    return idsArr.indexOf(item.from) > -1 || idsArr.indexOf(item.to) > -1;
                });
                // console.log(filteredEdges, "filtered edges")
                let leftOutIdsArr = [];
                for (let item of filteredEdges) {
                    if (leftOutIdsArr.indexOf(item.from) === -1) {
                        leftOutIdsArr.push(item.from);
                    }
                    if (leftOutIdsArr.indexOf(item.to) === -1) {
                        leftOutIdsArr.push(item.to);
                    }
                }
                // console.log(leftOutIdsArr, "leftOutIdsArr")
                filteredNodes = this.nodes.filter(function (item) {
                    return leftOutIdsArr.indexOf(item.id) > -1;
                })
                // console.log(filteredNodes, "refiltered nodes")
            } else {
                filteredNodes = this.nodes;
                filteredEdges = this.edges;
            }
            //console.log(filteredNodes);
            //console.log(filteredEdges);

            this.nodesDataset = new DataSet(filteredNodes);
            this.edgesDataset = new DataSet(filteredEdges);
            this.data = {
                nodes: this.nodesDataset,
                edges: this.edgesDataset
            };
            this.network = new Network(this.container.nativeElement, this.data, this.options);
            this.network.stabilize(100);
            // if we click on a node, we want to open it up!
            let self = this;
            this.network.on("selectNode", function (params) {
                if (params.nodes.length === 1 && self.network.isCluster(params.nodes[0])) {
                    self.network.openCluster(params.nodes[0]);
                }
            });
        }
    }

    public reset() {
        this.serviceName = 'ALL';
        this.enableClustering = false;
        this.reload('pod');
    }

    public service() {
        this.serviceName = 'ALL';
        this.enableClustering = true;
        this.reload('service');
    }

    public reload(nodeType) {
        this.clusterIndex = 0;
        this.clusters = [];
        this.lastClusterZoomLevel = 0;
        this.clusterFactor = 0.9;
        this.nodes = []
        this.edges = []
        this.nodesDataset = []
        this.edgesDataset = []

        this.NODE_STATUS = STATUS_WAIT;
        this.EDGE_STATUS = STATUS_WAIT;
        this.data = {};
        let nodeShape: string = 'dot';
        if (nodeType == 'service') {
            nodeShape = 'star';
        }
        this.options = {
            nodes: {
                shape: nodeShape,
                size: 16
            },
            physics: {
                enabled: false,
                /*forceAtlas2Based: {
                    gravitationalConstant: -26,
                    centralGravity: 0.005,
                    springLength: 230,
                    springConstant: 0.18
                },
                maxVelocity: 146,
                solver: 'forceAtlas2Based',
                timestep: 0.35,
                stabilization: { iterations: 150 }*/
            },
            layout: {
                improvedLayout: false
            }
        };
        this.network = {};
        this.loadApp(nodeType);
    }

    private loadApp(nodeType) {
        this.getNodes(nodeType);
        this.getEdges(nodeType);
    }

    ngOnInit() {
        //this.getServiceList();
        this.loadApp('pod');
    }

    public clusterByCid() {
        this.enableClustering = true;
        this.network.setData(this.data);
        let nodeServices = [];
        for (let item of this.nodes) {
            for (let cid of item.cid) {
                if (cid && nodeServices.indexOf(cid) === -1) {
                    nodeServices.push(cid);
                }
            }
        }
        /*for (let i = 0; i < this.nodes.length; i++) {
            nodeServices[i] = (this.nodes[i].cid);
        }*/
        let uniqServices = ([...nodeServices]);
        let clusterOptionsByData = new Array(uniqServices.length);
        for (let i = 0; i < uniqServices.length; i++) {
            clusterOptionsByData[i] = {
                joinCondition: function (childOptions) {
                    //return childOptions.cid == uniqServices[i];
                    return childOptions.cid && childOptions.cid.indexOf(uniqServices[i]) > -1;
                },
                clusterNodeProperties: { allowSingleNodeCluster: true, id: 'cidCluster' + i, borderWidth: 2, shape: 'star', label: uniqServices[i] }
            };
            //console.log(clusterOptionsByData[i]);
            this.network.cluster(clusterOptionsByData[i]);
        }
    }

}
