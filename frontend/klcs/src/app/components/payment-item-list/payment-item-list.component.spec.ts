import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PaymentItemListComponent } from './payment-item-list.component';

describe('PaymentItemListComponent', () => {
  let component: PaymentItemListComponent;
  let fixture: ComponentFixture<PaymentItemListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PaymentItemListComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PaymentItemListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
