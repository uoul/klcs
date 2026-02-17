import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateShopDialogComponent } from './create-shop-dialog.component';

describe('CreateShopDialogComponent', () => {
  let component: CreateShopDialogComponent;
  let fixture: ComponentFixture<CreateShopDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateShopDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateShopDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
