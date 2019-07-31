import { CompareModule } from './compare.module';

describe('CompareModule', () => {
  let compareModule: CompareModule;

  beforeEach(() => {
    compareModule = new CompareModule();
  });

  it('should create an instance', () => {
    expect(compareModule).toBeTruthy();
  });
});
