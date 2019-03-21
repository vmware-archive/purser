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
  public toBeDeletedGroup = "Custom Group";
  public groupCreation = 'wait';
  public groupDeletion = 'wait';

  public group: any;

  constructor(private router: Router, private logicalGroupService: LogicalGroupService) { }

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
      }, (err) => {
          console.log("failed to create group", err);
          this.groupCreation = 'fail';
      });
    this.isCreateGroup = false;
  }

  public deleteGroup() {
    console.log("deleting group ", this.toBeDeletedGroup)
    let observableEntity: Observable<any> = this.logicalGroupService.deleteCustomGroup(this.toBeDeletedGroup);
      observableEntity.subscribe((response) => {
          console.log("successfully deleted group");
          this.groupDeletion = 'success';
      }, (err) => {
          console.log("failed to delete group", err);
          this.groupDeletion = 'fail';
      });
    this.isDeleteGroup = false;
  }

  public setToBeDeletedGroup(grpName) {
    this.toBeDeletedGroup = grpName
  }

  public reset() {
    this.isCreateGroup = false;
    this.getLogicalGroupData();
    this.isDeleteGroup = false;
    this.toBeDeletedGroup = "Custom Group";
    this.group = null;
    this.groupCreation = 'wait';
    this.groupDeletion = 'wait';
  }

  ngOnInit() {
    this.reset()
  }

}
