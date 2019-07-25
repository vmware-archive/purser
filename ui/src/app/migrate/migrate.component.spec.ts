import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MigrateComponent } from './migrate.component';

describe('MigrateComponent', () => {
  let component: MigrateComponent;
  let fixture: ComponentFixture<MigrateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MigrateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MigrateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
