import { TestBed } from '@angular/core/testing';

import { ShopAdminApiService } from './shop-admin-api.service';

describe('ShopAdminApiService', () => {
  let service: ShopAdminApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ShopAdminApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
