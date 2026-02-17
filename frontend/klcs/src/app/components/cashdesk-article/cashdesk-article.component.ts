import { Component, input, output } from '@angular/core';
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'klcs-cashdesk-article',
  imports: [
    CommonModule,
  ],
  templateUrl: './cashdesk-article.component.html',
  styleUrl: './cashdesk-article.component.css'
})
export class CashdeskArticleComponent {
  article = input.required<Article>()
  articleClicked = output<Article>()
}
