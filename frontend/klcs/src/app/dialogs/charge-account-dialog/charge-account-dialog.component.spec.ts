import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChargeAccountDialogComponent } from './charge-account-dialog.component';

describe('ChargeAccountDialogComponent', () => {
  let component: ChargeAccountDialogComponent;
  let fixture: ComponentFixture<ChargeAccountDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChargeAccountDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ChargeAccountDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
