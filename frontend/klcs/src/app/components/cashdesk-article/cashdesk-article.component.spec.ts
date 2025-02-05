import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CashdeskArticleComponent } from './cashdesk-article.component';

describe('CashdeskArticleComponent', () => {
  let component: CashdeskArticleComponent;
  let fixture: ComponentFixture<CashdeskArticleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CashdeskArticleComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CashdeskArticleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
