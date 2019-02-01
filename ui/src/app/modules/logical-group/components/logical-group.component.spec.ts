import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LogicalGroupComponent } from './logical-group.component';

describe('LogicalGroupComponent', () => {
  let component: LogicalGroupComponent;
  let fixture: ComponentFixture<LogicalGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LogicalGroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LogicalGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
