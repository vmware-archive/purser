import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';
import { CompareService } from '../../services/compare.service';
import { Observable } from 'rxjs';
import { CloudDetails } from './cloud-details';
import { setDefaultService } from 'selenium-webdriver/opera';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss']
})
export class CompareCloudsComponent implements OnInit {

  regions :any;
  showCloud : boolean = false;
  detailsL = ["CPU", "Memory", "CPU Cost", "Memory Cost", "Total Cost"];
  showDetailsModal : boolean = false;
  showBtn : boolean = true;
  showBack : boolean = false;

  selectedRegions : any[] = Object.create(null);
  cloudDetails : any[] = [];
  cloudRegions : any[] = [];
  diffPercent : any[] = [];
  costDiff : any[] = [];

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
    var c;
    for(c = 0;c < this.cloudRegions.length; c++){
        this.selectedRegions[c] = "US-East-1";
    }
    this.cloudDetails = [
      {
        cloud : "AWS",
        cpu : 1,
        cpuCost : 100,
        memory : 20,
        memoryCost : 40,
        totalCost : 100,
        existingCost : 20
      },
      {
        cloud : "GCP",
        cpu : 1,
        cpuCost : 100,
        memory : 20,
        memoryCost : 40,
        totalCost : 100,
        existingCost : 100
      },
      {
        cloud : "PKS",
        cpu : 1,
        cpuCost : 100,
        memory : 20,
        memoryCost : 40,
        totalCost : 100,
        existingCost : 200
      },
      {
        cloud : "Azure",
        cpu : 1,
        cpuCost : 100,
        memory : 20,
        memoryCost : 40,
        totalCost : 100,
        existingCost : 120
      }
    ]
    this.cloudRegions = [
      {
        cloud : "Amazon Web Services",
        region : ["US-East-1", "US-West-2", "EU-West-1"]
      },
      {
        cloud : "Google Cloud Platform",
        region : ["US-East-1", "US-West-2", "EU-West-1"]      
      },
      {
        cloud : "Pivotal Container Service",
        region : ["US-East-1", "US-West-2", "EU-West-1"]
      },
      {
        cloud : "Microsoft Azure",
        region : ["US-East-1", "US-West-2", "EU-West-1"]      
      }
    ];
  }
  
  showClouds(){
    console.log("----selected values-----" + JSON.stringify(this.selectedRegions))
    for(var cd = 0; cd < this.cloudDetails.length; cd++){
      this.costDiff[cd] = this.cloudDetails[cd].totalCost - this.cloudDetails[cd].existingCost;
      if(this.costDiff[cd] < 0){

      }
    }

    this.showBtn = false;
    this.showCloud = true;
    this.showBack = true;
  }

  showDetails(){
    this.showDetailsModal = true;
  }
  back(){
    this.showBtn = true;
    this.showCloud = false;
    this.showBack = false;
  }
}