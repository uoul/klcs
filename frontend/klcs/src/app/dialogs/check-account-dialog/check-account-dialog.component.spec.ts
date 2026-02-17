import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckAccountDialogComponent } from './check-account-dialog.component';

describe('CheckAccountDialogComponent', () => {
  let component: CheckAccountDialogComponent;
  let fixture: ComponentFixture<CheckAccountDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CheckAccountDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CheckAccountDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
