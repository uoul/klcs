import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CashdeskComponent } from './cashdesk.component';

describe('CashdeskComponent', () => {
  let component: CashdeskComponent;
  let fixture: ComponentFixture<CashdeskComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CashdeskComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CashdeskComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
