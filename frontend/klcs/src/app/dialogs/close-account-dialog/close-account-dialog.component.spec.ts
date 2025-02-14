import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CloseAccountDialogComponent } from './close-account-dialog.component';

describe('CloseAccountDialogComponent', () => {
  let component: CloseAccountDialogComponent;
  let fixture: ComponentFixture<CloseAccountDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CloseAccountDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CloseAccountDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
