import { Component, OnInit } from '@angular/core';
import * as $ from 'jquery';
import { Observable } from 'rxjs';
import { PlanInfraService } from '../../services/plan-infra.service';

@Component({
  selector: 'app-plan-infra',
  templateUrl: './plan-infra.component.html',
  styleUrls: ['./plan-infra.component.scss'],
  providers:[PlanInfraService]
})
export class PlanInfraComponent implements OnInit {

  fileToUpload: File = null;

  cloudRegions = [
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

  constructor(private planInfraService : PlanInfraService) { }

  ngOnInit() {
  }

  handleFileInput(files: FileList) {
    this.fileToUpload = files.item(0);
  }
  
  uploadFileToActivity() {
    this.planInfraService.postFile(this.fileToUpload).subscribe(data => {
      // do something, if upload success
      }, error => {
        console.log(error);
      });
  }


}
