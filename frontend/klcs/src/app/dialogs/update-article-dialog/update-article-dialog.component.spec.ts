import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UpdateArticleDialogComponent } from './update-article-dialog.component';

describe('UpdateArticleDialogComponent', () => {
  let component: UpdateArticleDialogComponent;
  let fixture: ComponentFixture<UpdateArticleDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UpdateArticleDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UpdateArticleDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
