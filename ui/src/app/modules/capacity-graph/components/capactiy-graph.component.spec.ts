import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CapactiyGraphComponent } from './capactiy-graph.component';

describe('CapactiyGraphComponent', () => {
  let component: CapactiyGraphComponent;
  let fixture: ComponentFixture<CapactiyGraphComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CapactiyGraphComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CapactiyGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
