import { Component, OnInit } from '@angular/core';
import { AppProfile } from '../app-profile';
import { AppProfileService } from '../../services/app-profile.service';
import { CloudCost } from './cloud-cost';

@Component({
  selector: 'app-run-scenario',
  templateUrl: './run-scenario.component.html',
  styleUrls: ['./run-scenario.component.scss']
})
export class RunScenarioComponent implements OnInit {

  scenarioName : string;
  appProfile : AppProfile;
  images : string[] = ["aws.jpg", "azure.png", "gcp.jpg", "ibm.png", "vmc.png"];

  cloudCost : CloudCost[] = [
    {
      cloud : "AWS",
      cost : 120
    },
    {
      cloud : "Azure",
      cost : 150      
    },
    {
      cloud : "GCP",
      cost : 100     
    },
    {
      cloud : "IBM",
      cost : 200      
    },
    {
      cloud : "VMC",
      cost : 190      
    }  
  ]

  constructor(private appProService : AppProfileService) { }

  ngOnInit() {
    this.scenarioName = this.appProService.scenarioName;
    this.appProfile = this.appProService.appProfile;
  }

}