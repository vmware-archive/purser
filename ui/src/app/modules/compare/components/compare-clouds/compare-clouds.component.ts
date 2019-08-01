import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';
import { CompareService } from '../../services/compare.service';
import { Observable } from 'rxjs';
import { CloudDetails } from './cloud-details';
import { setDefaultService } from 'selenium-webdriver/opera';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss'],
  providers:[CompareService]
})
export class CompareCloudsComponent implements OnInit {

  regions :any;
  showCloud : boolean = false;
  showDetailsModal : boolean = false;
  showBtn : boolean = true;
  showBack : boolean = false;
  nodes : any[] = [];
  cloudDetails : any[] = [];
  cloudRegions : any[] = [];
  diffPercent : any[] = [];
  costDiff : any[] = [];

  detailsResponse : any[] = [];

  sendCloudRegion : any[] = [];

  images = ["awst.png", "gcpt.png", "pkst.png", "azuret.png"];

  myStyles = [{
    'background-color': '#FEF3B5',
    },
    {
      'background-color': '#E1F1F6',
    },
    {
      'background-color': '#DFF0D0',
    },
    {
      'background-color': '#F5DBD9',
    }
  ]

  cardColors = [
    "'backgroudColor' : '#E1F1F6'",
    "'backgroudColor' : '#FEF3B5'",
    "'backgroudColor' : '#F5DBD9'",
    "'backgroudColor' : '#DFF0D0'",
  ]
  
  constructor(private compareService : CompareService) { }

  ngOnInit() {

    this.setDefault();

    this.regions = this.compareService.getRegions().subscribe(response => {
      console.log("Regions for clouds" + response);
    });

  }

  setDefault(){
    this.sendCloudRegion = [];

    this.cloudRegions = [
      {
        cloud : "Amazon Web Services",
        region : ["US-East-1", "US-West-2", "EU-West-1"],
        selectedRegion : "US-East-1"
      },
      {
        cloud : "Google Cloud Platform",
        region : ["US-East-1", "US-West-2", "EU-West-1"],
        selectedRegion : "US-East-1"      
      },
      {
        cloud : "Pivotal Container Service",
        region : ["US-East-1", "US-West-2", "EU-West-1"],
        selectedRegion : "US-East-1"
      },
      {
        cloud : "Microsoft Azure",
        region : ["US-East-1", "US-West-2", "EU-West-1"] ,
        selectedRegion : "US-East-1"     
      }
    ];
    this.cloudDetails = 

    [
    
      {
    
        "cloud": "aws",
    
        "existingCost": 1000,
    
        "totalCost": 48.96,
    
        "cpuCost": 34.56,
    
        "memoryCost": 14.4,
    
        "cpu": 2,
    
        "memory": 2,
    
        "nodes": [
    
          {
    
            "instanceType": "t3.small",
    
            "os": "Windows",
    
            "totalNodeCost": 48.96,
    
            "cpuNodeCost": 34.56,
    
            "memoryNodeCost": 14.4,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          },
          {
    
            "instanceType": "t3.small",
    
            "os": "Windows",
    
            "totalNodeCost": 48.96,
    
            "cpuNodeCost": 34.56,
    
            "memoryNodeCost": 14.4,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          },
          {
    
            "instanceType": "t3.small",
    
            "os": "Windows",
    
            "totalNodeCost": 48.96,
    
            "cpuNodeCost": 34.56,
    
            "memoryNodeCost": 14.4,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          },
          {
    
            "instanceType": "t3.small",
    
            "os": "Windows",
    
            "totalNodeCost": 48.96,
    
            "cpuNodeCost": 34.56,
    
            "memoryNodeCost": 14.4,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          }
    
        ]
    
      },
    
      {
    
        "cloud": "gcp",
    
        "existingCost": 1000,
    
        "totalCost": 51.621120000000005,
    
        "cpuCost": 45.51984,
    
        "memoryCost": 6.10128,
    
        "cpu": 2,
    
        "memory": 2,
    
        "nodes": [
    
          {
    
            "instanceType": "n1-standard",
    
            "os": "linux",
    
            "totalNodeCost": 51.62112,
    
            "cpuNodeCost": 45.51984,
    
            "memoryNodeCost": 6.10128,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          }
    
        ]
    
      },
    
      {
    
        "cloud": "pks",
    
        "existingCost": 1000,
    
        "totalCost": 74.448,
    
        "cpuCost": 69.264,
    
        "memoryCost": 5.184,
    
        "cpu": 2,
    
        "memory": 2,
    
        "nodes": [
    
          {
    
            "instanceType": "PKS-US-East-1",
    
            "os": "linux",
    
            "totalNodeCost": 74.448,
    
            "cpuNodeCost": 69.264,
    
            "memoryNodeCost": 5.184,
    
            "cpuNode": 2,
    
            "memoryNode": 2
    
          }
    
        ]
    
      },
    
      {
    
        "cloud": "azure",
    
        "existingCost": 1000,
    
        "totalCost": 59.760000000000005,
    
        "cpuCost": 34.56,
    
        "memoryCost": 25.200000000000003,
    
        "cpu": 2,
    
        "memory": 3.5,
    
        "nodes": [
    
          {
    
            "instanceType": "Basic_A2",
    
            "os": "windows",
    
            "totalNodeCost": 59.760000000000005,
    
            "cpuNodeCost": 34.56,
    
            "memoryNodeCost": 25.200000000000003,
    
            "cpuNode": 2,
    
            "memoryNode": 3.5
    
          }
    
        ]
    
      }
    
    ]
        
    /*
    var c;
    for(c = 0;c < this.cloudRegions.length; c++){
        this.selectedRegions[c] = "US-East-1";
    }
    */  
  }
  
  showClouds(){

    this.showBtn = false;
    this.showCloud = true;
    this.showBack = true;
    
    for(let cd of this.cloudDetails){
      cd.costDiff = (cd.totalCost - cd.existingCost).toFixed(2);
      cd.costPercent = ((cd.costDiff / cd.totalCost) * 100).toFixed(2);
    }
    /*
    for(var c in this.cloudRegions ){
      this.sendCloudRegion.push({
        'cloud': this.cloudRegions[c].cloud,
        'region': this.selectedRegions[c]
      });
    }
    */
    this.compareService.regions = this.cloudRegions;
    this.compareService.sendCloudRegion(this.sendCloudRegion).subscribe(data => {
        console.log(data);
        this.compareService.cloudDetails = data;
    });
  }

  showDetails(cloud){
    this.nodes = cloud.nodes;
    this.showDetailsModal = true;
  }

  back(){
    this.showBtn = true;
    this.showCloud = false;
    this.showBack = false;
    this.setDefault();
  }
}