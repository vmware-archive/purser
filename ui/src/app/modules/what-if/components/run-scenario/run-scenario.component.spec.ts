import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RunScenarioComponent } from './run-scenario.component';

describe('RunScenarioComponent', () => {
  let component: RunScenarioComponent;
  let fixture: ComponentFixture<RunScenarioComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RunScenarioComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RunScenarioComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
