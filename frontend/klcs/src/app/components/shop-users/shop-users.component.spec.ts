import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShopUsersComponent } from './shop-users.component';

describe('ShopUsersComponent', () => {
  let component: ShopUsersComponent;
  let fixture: ComponentFixture<ShopUsersComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShopUsersComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShopUsersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
