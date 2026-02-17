import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReadQrDialogComponent } from './read-qr-dialog.component';

describe('ReadQrDialogComponent', () => {
  let component: ReadQrDialogComponent;
  let fixture: ComponentFixture<ReadQrDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReadQrDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReadQrDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
