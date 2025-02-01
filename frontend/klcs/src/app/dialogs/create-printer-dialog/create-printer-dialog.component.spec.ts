import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePrinterDialogComponent } from './create-printer-dialog.component';

describe('CreatePrinterDialogComponent', () => {
  let component: CreatePrinterDialogComponent;
  let fixture: ComponentFixture<CreatePrinterDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreatePrinterDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreatePrinterDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
