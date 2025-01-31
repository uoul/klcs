import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShopPrintersComponent } from './shop-printers.component';

describe('ShopPrintersComponent', () => {
  let component: ShopPrintersComponent;
  let fixture: ComponentFixture<ShopPrintersComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShopPrintersComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShopPrintersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
