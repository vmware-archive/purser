import { PlanInfraModule } from './plan-infra.module';

describe('PlanInfraModule', () => {
  let planInfraModule: PlanInfraModule;

  beforeEach(() => {
    planInfraModule = new PlanInfraModule();
  });

  it('should create an instance', () => {
    expect(planInfraModule).toBeTruthy();
  });
});
