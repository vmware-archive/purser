import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';
import { CompareService } from '../../services/compare.service';
import { Observable } from 'rxjs';
import { CloudDetails } from './cloud-details';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss']
})
export class CompareCloudsComponent implements OnInit {
  regions :any;
  showCloud : boolean = false;
  detailsL = ["CPU", "Memory", "CPU Cost", "Memory Cost", "Total Cost"];
  basic : boolean = false;
  showBtn : boolean = true;

  selectedRegions : any[] = Object.create(null);

  cloudDetails = [
    {
      cloud : "AWS",
      cpu : 1,
      cpuCost : 100,
      memory : 20,
      memoryCost : 40,
      totalCost : 100
    },
    {
      cloud : "AWS",
      cpu : 1,
      cpuCost : 100,
      memory : 20,
      memoryCost : 40,
      totalCost : 100
    },
    {
      cloud : "AWS",
      cpu : 1,
      cpuCost : 100,
      memory : 20,
      memoryCost : 40,
      totalCost : 100
    },
    {
      cloud : "AWS",
      cpu : 1,
      cpuCost : 100,
      memory : 20,
      memoryCost : 40,
      totalCost : 100
    }
  ]

  cloudRegions : CloudRegion[] = [
    {
      cloud : "AWS",
      region : ["US-East-1", "US-West-2", "EU-West-1"]
    },
    {
      cloud : "GCP",
      region : ["US-East-1", "US-West-2", "EU-West-1"]      
    },
    {
      cloud : "PKS",
      region : ["US-East-1", "US-West-2", "EU-West-1"]
    },
    {
      cloud : "Azure",
      region : ["US-East-1", "US-West-2", "EU-West-1"]      
    }
  ];

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

    var c;
    for(c = 0;c < this.cloudRegions.length; c++){
        this.selectedRegions[c] = "US-East-1";
    }

    this.regions = this.compareService.getRegions().subscribe(response => {
      console.log("Regions for clouds" + response);
    });

  }
  
  showClouds(){
    console.log("----selected values-----" + JSON.stringify(this.selectedRegions))
    this.showBtn = false;
    this.showCloud = true;
  }

  showDetails(){
    this.basic = true;
  }
}