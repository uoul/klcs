import { Component, input, signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';
import { CashdeskArticleComponent } from "../cashdesk-article/cashdesk-article.component";
import { CommonModule } from '@angular/common';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { ShoppingCartComponent } from "../shopping-cart/shopping-cart.component";
import { EditCartDialogComponent } from "../../dialogs/edit-cart-dialog/edit-cart-dialog.component";
import { CheckoutDialogComponent } from "../../dialogs/checkout-dialog/checkout-dialog.component";
import { SellerApiService } from '../../services/seller-api/seller-api.service';

@Component({
  selector: 'klcs-cashdesk',
  imports: [
    CommonModule,
    CashdeskArticleComponent,
    ShoppingCartComponent,
    EditCartDialogComponent,
    CheckoutDialogComponent
],
  templateUrl: './cashdesk.component.html',
  styleUrl: './cashdesk.component.css'
})
export class CashdeskComponent {

  protected readonly EDIT_CART_DIALOG_ID = "edit-cart-dialog"
  protected readonly CHECKOUT_DIALOG_ID = "checkout-dialog"

  checkoutCard: WritableSignal<boolean> = signal(false)

  constructor(
    protected shoppingCart: ShoppingCartService,
    protected sellerApi: SellerApiService,
  ){}

  addToShoppingCart(article: Article) {
    this.shoppingCart.addArticle(article)
  }

  showShoppingCartDialog() {
    const dialog = document.getElementById(this.EDIT_CART_DIALOG_ID) as HTMLDialogElement
    dialog.showModal()
  }

  showCheckout(withCard: boolean){
    this.checkoutCard.set(withCard)
    const dialog = document.getElementById(this.CHECKOUT_DIALOG_ID) as HTMLDialogElement
    dialog.showModal()
  }

  checkCategoryEmpty(articles: Article[]): boolean {
    for(let article of articles){
      if(article.StockAmount === null || article.StockAmount > 0){
        return false
      }
    }
    return true
  }
}
