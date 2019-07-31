import { Component, OnInit } from '@angular/core';
import { AppProfileService } from '../../services/app-profile.service';
import {  Router } from '@angular/router';
import { AppProfile } from '../app-profile';

@Component({
  selector: 'app-what-if',
  templateUrl: './what-if.component.html',
  styleUrls: ['./what-if.component.scss'],
  providers: [AppProfileService]
})
export class WhatIfComponent implements OnInit {

  addWorkload : string = "Add Workload"
  appProfileTitle : string = "Application Profile"
  scenarioName : string;
  appProfile : any = {};
  
  constructor(private appProService : AppProfileService, private router : Router) { 
  }

  runScenario(): void{

    //var sendProfile = JSON.stringify(this.appProfile)

    this.appProService.scenarioName = this.scenarioName;
    this.appProService.appProfile = this.appProfile;
    console.log("-----appProfile-----" + JSON.stringify(this.appProfile))


    this.appProService.submitSpec(this.appProfile).subscribe((response) =>{
      console.log("response" + response);
    })

    this.router.navigate(['/runscene']); 
  }

  ngOnInit() {
    console.log(this.appProfile);
    this.appProfile.startDate = new Date();
  }

   preventNonNumericalInput(e) {
    e = e || window.event;
    var charCode = (typeof e.which == "undefined") ? e.keyCode : e.which;
    var charStr = String.fromCharCode(charCode);
  
    if (!charStr.match(/^[0-9]+$/))
      e.preventDefault();
  }

}