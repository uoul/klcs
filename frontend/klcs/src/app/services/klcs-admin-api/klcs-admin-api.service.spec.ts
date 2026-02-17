import { TestBed } from '@angular/core/testing';

import { KlcsAdminApiService } from './klcs-admin-api.service';

describe('KlcsAdminApiService', () => {
  let service: KlcsAdminApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KlcsAdminApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
