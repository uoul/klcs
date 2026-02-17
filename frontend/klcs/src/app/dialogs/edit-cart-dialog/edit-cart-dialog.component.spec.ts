import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditCartDialogComponent } from './edit-cart-dialog.component';

describe('EditCartDialogComponent', () => {
  let component: EditCartDialogComponent;
  let fixture: ComponentFixture<EditCartDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditCartDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditCartDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
