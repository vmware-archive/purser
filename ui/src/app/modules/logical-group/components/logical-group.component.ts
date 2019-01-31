import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { LogicalGroupService } from '../services/logical-group.service';

const STATUS_WAIT = 'WAIT',
    STATUS_READY = 'READY',
    STATUS_NODATA = 'NO_DATA';

@Component({
  selector: 'app-logical-group',
  templateUrl: './logical-group.component.html',
  styleUrls: ['./logical-group.component.css']
})
export class LogicalGroupComponent implements OnInit {

  title = 'my-app';
  groups: Object[];
  public GROUP_STATUS = STATUS_WAIT;

  constructor(private router: Router, private capacityGraphService: LogicalGroupService) { }

    private getLogicalGroupData() {
        let observableEntity: Observable<any> = this.capacityGraphService.getLogicalGroupData();
        this.GROUP_STATUS = STATUS_WAIT;
        observableEntity.subscribe((response) => {
            if (!response) {
              console.log("empty response")
                return;
            }
            this.groups = JSON.parse(JSON.stringify(response));
            console.log(this.groups)
        }, (err) => {
            this.GROUP_STATUS = STATUS_NODATA;
        });
    }


  public sortName() {
    this.groups.sort((obj1: any, obj2: any) => {
      if (obj1.name > obj2.name) {
        return 1;
      }

      if (obj1.name < obj2.name) {
        return -1;
      }
      return 0;
    })
  }

  public revSortName() {
    this.groups.sort((obj1: any, obj2: any) => {
      if (obj1.name < obj2.name) {
        return 1;
      }

      if (obj1.name > obj2.name) {
        return -1;
      }
      return 0;
    })
  }

  public sortPodsCount() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj1.podsCount - obj2.podsCount;
    })
  }

  public revSortPodsCount() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj2.podsCount - obj1.podsCount;
    })
  }

  public sortCPU() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj1.cpu - obj2.cpu;
    })
  }

  public revSortCPU() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj2.cpu - obj1.cpu;
    })
  }

  public sortMemory() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj1.memory - obj2.memory;
    })
  }

  public revSortMemory() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj2.memory - obj1.memory;
    })
  }

  public sortStorage() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj1.storage - obj2.storage;
    })
  }

  public revSortStorage() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj2.storage - obj1.storage;
    })
  }

  public sortMTDCost() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj1.mtdCost - obj2.mtdCost;
    })
  }

  public revMTDCost() {
    this.groups.sort((obj1: any, obj2: any) => {
      return obj2.mtdCost - obj1.mtdCost;
    })
  }

  ngOnInit() {
    this.getLogicalGroupData()
  }

}
