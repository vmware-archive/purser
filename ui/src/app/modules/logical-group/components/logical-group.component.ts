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

  ngOnInit() {
    this.getLogicalGroupData()
  }

}
