import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UpdateShopDialogComponent } from './update-shop-dialog.component';

describe('UpdateShopDialogComponent', () => {
  let component: UpdateShopDialogComponent;
  let fixture: ComponentFixture<UpdateShopDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UpdateShopDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UpdateShopDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
