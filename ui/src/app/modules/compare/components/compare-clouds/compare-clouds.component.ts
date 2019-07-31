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
  detailsL = ["CPU", "Memory", "CPU Cost", "Memory Cost", "Total Cost"];
  showDetailsModal : boolean = false;
  showBtn : boolean = true;
  showBack : boolean = false;

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
    /*
    var c;
    for(c = 0;c < this.cloudRegions.length; c++){
        this.selectedRegions[c] = "US-East-1";
    }
    */
    console.log("------default-------" + JSON.stringify(this.cloudRegions))  
  }
  
  showClouds(){

    this.showBtn = false;
    this.showCloud = true;
    this.showBack = true;
    
    for(let cd of this.cloudDetails){
      cd.costDiff = cd.totalCost - cd.existingCost;
      console.log("-----cloud details---" + JSON.stringify(cd));
    }
    /*
    for(var c in this.cloudRegions ){
      this.sendCloudRegion.push({
        'cloud': this.cloudRegions[c].cloud,
        'region': this.selectedRegions[c]
      });
    }
    */
    /*
    this.compareService.sendCloudRegion(this.sendCloudRegion).subscribe(data => {
        console.log(data);
    });
    */ 
    console.log("--------post data-------" + JSON.stringify(this.cloudRegions));
  }

  showDetails(){
    this.showDetailsModal = true;
  }
  back(){
    this.showBtn = true;
    this.showCloud = false;
    this.showBack = false;

    this.setDefault();
  }
}