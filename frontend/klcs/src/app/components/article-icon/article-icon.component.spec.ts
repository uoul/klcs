import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ArticleIconComponent } from './article-icon.component';

describe('ArticleIconComponent', () => {
  let component: ArticleIconComponent;
  let fixture: ComponentFixture<ArticleIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ArticleIconComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ArticleIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
