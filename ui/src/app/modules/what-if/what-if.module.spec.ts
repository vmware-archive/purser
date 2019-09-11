import { WhatIfModule } from './what-if.module';

describe('WhatIfModule', () => {
  let whatIfModule: WhatIfModule;

  beforeEach(() => {
    whatIfModule = new WhatIfModule();
  });

  it('should create an instance', () => {
    expect(whatIfModule).toBeTruthy();
  });
});
