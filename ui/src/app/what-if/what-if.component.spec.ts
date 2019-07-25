import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WhatIfComponent } from './what-if.component';

describe('WhatIfComponent', () => {
  let component: WhatIfComponent;
  let fixture: ComponentFixture<WhatIfComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WhatIfComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WhatIfComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
