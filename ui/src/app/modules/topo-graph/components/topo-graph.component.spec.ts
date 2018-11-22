import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TopoGraphComponent } from './topo-graph.component';

describe('TopoGraphComponent', () => {
  let component: TopoGraphComponent;
  let fixture: ComponentFixture<TopoGraphComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TopoGraphComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TopoGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
