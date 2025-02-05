import { Component, input } from '@angular/core';
import { Article } from '../../domain/Article';
import { CashdeskArticleComponent } from "../cashdesk-article/cashdesk-article.component";
import { CommonModule } from '@angular/common';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { ShoppingCartComponent } from "../shopping-cart/shopping-cart.component";

@Component({
  selector: 'klcs-cashdesk',
  imports: [
    CommonModule,
    CashdeskArticleComponent,
    ShoppingCartComponent
],
  templateUrl: './cashdesk.component.html',
  styleUrl: './cashdesk.component.css'
})
export class CashdeskComponent {
  categories = input.required<Map<string, Article[]>>()

  constructor(
    private shoppingCart: ShoppingCartService,
  ){}

  addToShoppingCart(article: Article) {
    this.shoppingCart.addArticle(article)
  }
}
