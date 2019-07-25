import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-what-if',
  templateUrl: './what-if.component.html',
  styleUrls: ['./what-if.component.css']
})
export class WhatIfComponent implements OnInit {

  addWorkload : string = "Add Workload"
  appProfile : string = "Application Profile"
  scenarioName : string;
  profile_labels = ["CPU", "Memory", "Disk Space", "Annual Project Growth", "Number of VMs"];
  profile_units = ["vCPU", "GB", "GB", "%", ""];
  profile : any[] = [];
  profile_values : number[] = [];
  startDate : any;
  endDate : any;
  
  constructor() { 
  }

  runScenario(): void{
    console.log("------given values------" + this.profile_values);
  }

  ngOnInit() {
    for(var i = 0; i <this.profile_labels.length; i++){
        this.profile.push({
          label : this.profile_labels[i],
          unit : this.profile_units[i]
        });
    }
  }

}
