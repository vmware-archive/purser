import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss']
})
export class CompareCloudsComponent implements OnInit {

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
    "'backgroudColor' : '#DFF0D0'",
    "'backgroudColor' : '#F5DBD9'",
    "'backgroudColor' : '#FEF3B5'",
  ]
  
  constructor() { }

  ngOnInit() {
  }

}
