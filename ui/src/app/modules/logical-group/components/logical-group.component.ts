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
  public isCreateGroup = false;
  public isDeleteGroup = false;
  public isShowGroupDetails = false;
  public toBeDeletedGroup = "Custom Group";
  public groupToFocus: any;
  public groupCreation = 'wait';
  public groupDeletion = 'wait';
  public creationError = null;
  public deletionError = null;

  public isShowMTD = false;
  public isShowProjected = false;
  public donutOptions = {};
  public donutData = {"data": []};
  public group: any;
  public costRatio = 100;

  public isViewCurrentMonth = true;
  public isViewLastMonth = false;
  public isViewLastThreeMonth = false;
  public isViewProjected = false;

  constructor(private router: Router, private logicalGroupService: LogicalGroupService) {
  }

  public getCurrentMonthStatus() {
    if (this.isViewCurrentMonth) {
      return 'label-info'
    } else {
      return 'label-purple'
    }
  }

  public getLastMonthStatus() {
    if (this.isViewLastMonth) {
      return 'label-info'
    } else {
      return 'label-purple'
    }
  }

  public getLast3MonthStatus() {
    if (this.isViewLastThreeMonth) {
      return 'label-info'
    } else {
      return 'label-purple'
    }
  }

  public getProjectedMonthStatus() {
    if (this.isViewProjected) {
      return 'label-info'
    } else {
      return 'label-purple'
    }
  }

  public changeCurrentMonthView() {
    this.isViewCurrentMonth = !this.isViewCurrentMonth;
  }

  public changeLastMonthView() {
    this.isViewLastMonth = !this.isViewLastMonth;
  }

  public changeLast3MonthView() {
    this.isViewLastThreeMonth = !this.isViewLastThreeMonth;
  }

  public changeProjectedMonthView() {
    this.isViewProjected = !this.isViewProjected;
  }


  private getLogicalGroupData() {
      let observableEntity: Observable<any> = this.logicalGroupService.getLogicalGroupData();
      this.GROUP_STATUS = STATUS_WAIT;
      observableEntity.subscribe((response) => {
          if (!response) {
            console.log("empty response")
              return;
          }
          this.groups = JSON.parse(JSON.stringify(response));
      }, (err) => {
          this.GROUP_STATUS = STATUS_NODATA;
      });
  }

  public fillGroupData() {
    this.isCreateGroup = true;
    this.group = null;
  }

  public deleteGroupData() {
    this.toBeDeletedGroup = "Custom Group";
    this.isDeleteGroup = true;
  }

  public createGroup() {
    let observableEntity: Observable<any> = this.logicalGroupService.createCustomGroup(this.group);
      observableEntity.subscribe((response) => {
          console.log("successfully created group");
          this.groupCreation = 'success';
          setTimeout(() => this.ngOnInit(), 6000);
      }, (err) => {
          console.log("failed to create group", err);
          this.groupCreation = 'fail';
          this.creationError = err["error"];
          ;
      });
    this.isCreateGroup = false;
  }

  public deleteGroup() {
    console.log("deleting group ", this.toBeDeletedGroup)
    let observableEntity: Observable<any> = this.logicalGroupService.deleteCustomGroup(this.toBeDeletedGroup);
      observableEntity.subscribe((response) => {
          console.log("successfully deleted group");
          this.groupDeletion = 'success';
          setTimeout(() => this.ngOnInit(), 6000);
      }, (err) => {
          console.log("failed to delete group", err);
          this.groupDeletion = 'fail';
          this.deletionError = err["error"];
      });
    this.isDeleteGroup = false;
  }

  public setToBeDeletedGroup(grpName) {
    this.toBeDeletedGroup = grpName;
    this.isDeleteGroup = true;
  }

  public showGroupDetails(group) {
    console.log("group: ", group);
    this.groupToFocus = group;
    this.isShowGroupDetails = true;
    this.costRatio = Math.round(this.groupToFocus.mtdCost * 100 / this.groupToFocus.projectedCost);
  }

  public reset() {
    this.isCreateGroup = false;
    this.getLogicalGroupData();
    this.isDeleteGroup = false;
    this.isShowGroupDetails = false;
    this.toBeDeletedGroup = "Custom Group";
    this.group = null;
    this.groupCreation = 'wait';
    this.groupDeletion = 'wait';
  }

  public showMTD() {
    this.isShowMTD = true;
    this.isShowProjected = false;
    this.donutData = {
      "data": [
        ['CPU', this.groupToFocus.mtdCPUCost],
        ['Memory', this.groupToFocus.mtdMemoryCost],
        ['Storage', this.groupToFocus.mtdStorageCost]
      ]
    };

    this.donutOptions = {
      title: 'Total MTD Cost for ' + this.groupToFocus.name + ': ' + this.groupToFocus.mtdCost.toFixed(2),
      pieHole: 0.3,
      pieSliceText: 'value-and-percentage',
      width: 750,
      height: 400,
      chartArea: {
        left: "10%",
        top: "10%",
        height: "80%",
        width: "80%"
      }
    };
  }

  public showProjected() {
    this.isShowProjected = true;
    this.isShowMTD = false;
    this.donutData = {
      "data": [
        ['CPU', this.groupToFocus.projectedCPUCost],
        ['Memory', this.groupToFocus.projectedMemoryCost],
        ['Storage', this.groupToFocus.projectedStorageCost]
      ]
    };

    this.donutOptions = {
      title: 'Total Projected Cost for ' + this.groupToFocus.name + ': ' + this.groupToFocus.projectedCost.toFixed(2),
      pieHole: 0.3,
      pieSliceText: 'value-and-percentage',
      width: 750,
      height: 400,
      chartArea: {
        left: "10%",
        top: "10%",
        height: "80%",
        width: "80%"
      }
    };
  }

  ngOnInit() {
    this.reset();
  }

}
