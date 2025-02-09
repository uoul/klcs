import { Component, input } from '@angular/core';
import { Article } from '../../domain/Article';
import { CashdeskArticleComponent } from "../cashdesk-article/cashdesk-article.component";
import { CommonModule } from '@angular/common';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { ShoppingCartComponent } from "../shopping-cart/shopping-cart.component";
import { EditCartDialogComponent } from "../../dialogs/edit-cart-dialog/edit-cart-dialog.component";

@Component({
  selector: 'klcs-cashdesk',
  imports: [
    CommonModule,
    CashdeskArticleComponent,
    ShoppingCartComponent,
    EditCartDialogComponent
],
  templateUrl: './cashdesk.component.html',
  styleUrl: './cashdesk.component.css'
})
export class CashdeskComponent {
  categories = input.required<Map<string, Article[]>>()

  protected readonly EDIT_CART_DIALOG_ID = "edit-cart-dialog"

  constructor(
    protected shoppingCart: ShoppingCartService,
  ){}

  addToShoppingCart(article: Article) {
    this.shoppingCart.addArticle(article)
  }

  showShoppingCartDialog() {
    const dialog = document.getElementById(this.EDIT_CART_DIALOG_ID) as HTMLDialogElement
    dialog.showModal()
  }
}
