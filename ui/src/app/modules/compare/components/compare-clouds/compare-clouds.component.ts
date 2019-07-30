import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';
import { CompareService } from '../../services/compare.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss']
})
export class CompareCloudsComponent implements OnInit {
  regions :any;
  showCloud : boolean = false;

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

  images = ["aws.jpg", "gcp.jpg", "pks.jpg", "azure.png"];

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
    this.regions = this.compareService.getRegions().subscribe(response => {
      console.log("Regions for clouds" + this.regions);
    });
  }

  showClouds(){
    this.showCloud = true;
  }
}
