import { Component, OnInit } from '@angular/core';
import { CloudRegion } from './cloud-region';

@Component({
  selector: 'app-compare-clouds',
  templateUrl: './compare-clouds.component.html',
  styleUrls: ['./compare-clouds.component.scss']
})
export class CompareCloudsComponent implements OnInit {

  cloudRegions : CloudRegion[];
  
  constructor() { }

  ngOnInit() {
  }

}
