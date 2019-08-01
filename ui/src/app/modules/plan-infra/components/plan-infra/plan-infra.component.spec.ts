import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PlanInfraComponent } from './plan-infra.component';

describe('PlanInfraComponent', () => {
  let component: PlanInfraComponent;
  let fixture: ComponentFixture<PlanInfraComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PlanInfraComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PlanInfraComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
