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
  Arr = Array;
  public num:number = 1;
  private old:number = 0;
  public exprCount = [0];
  public exprStartIndices = [1];
  public toBeDeletedGroup = "Custom Group";

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
      }, (err) => {
          this.GROUP_STATUS = STATUS_NODATA;
      });
  }

  public fillGroupData() {
    this.isCreateGroup = true;
    this.num = 1;
    this.exprCount = [0];
  }

  public deleteGroupData() {
    this.toBeDeletedGroup = "Custom Group";
    this.isDeleteGroup = true;
  }

  public createGroup() {
    this.isCreateGroup = false;
  }

  public deleteGroup() {
    this.isDeleteGroup = false;
  }

  public increaseNum() {
    this.num++;
  }

  public increaseExpression() {
    this.exprCount.push(this.num - this.old);
    this.old = this.num;
    this.num++;
    this.exprStartIndices.push(this.num);
  }

  public getExprOfNum(n) {
    let exprNum = 0;
    for (let startIndex of this.exprStartIndices) {
      if (n < startIndex) {
        return exprNum;
      } else {
        exprNum++;
      }
    }
    return exprNum;
  }

  public setToBeDeletedGroup(grpName) {
    this.toBeDeletedGroup = grpName
  }

  ngOnInit() {
    this.isCreateGroup = false;
    this.num = 1;
    this.getLogicalGroupData();
    this.exprCount = [0];
    this.old = 0;
    this.isDeleteGroup = false;
    this.toBeDeletedGroup = "Custom Group";
  }

}
