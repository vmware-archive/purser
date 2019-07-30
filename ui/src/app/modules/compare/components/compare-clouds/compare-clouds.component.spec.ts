import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompareCloudsComponent } from './compare-clouds.component';

describe('CompareCloudsComponent', () => {
  let component: CompareCloudsComponent;
  let fixture: ComponentFixture<CompareCloudsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompareCloudsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompareCloudsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
