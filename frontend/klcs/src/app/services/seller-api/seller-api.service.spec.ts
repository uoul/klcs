import { TestBed } from '@angular/core/testing';

import { SellerApiService } from './seller-api.service';

describe('SellerApiService', () => {
  let service: SellerApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SellerApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
