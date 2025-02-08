import { Component, OnInit, signal, Signal, WritableSignal } from '@angular/core';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { OrderPosition } from '../../domain/OrderPosition';
import { ShoppingCartItemComponent } from "../shopping-cart-item/shopping-cart-item.component";
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';
import { CheckoutDialogComponent } from "../../dialogs/checkout-dialog/checkout-dialog.component";

@Component({
  selector: 'klcs-shopping-cart',
  imports: [
    CommonModule,
    ShoppingCartItemComponent,
    CheckoutDialogComponent
],
  templateUrl: './shopping-cart.component.html',
  styleUrl: './shopping-cart.component.css'
})
export class ShoppingCartComponent implements OnInit {
  
  readonly CHECKOUT_DIALOG_ID = "checkout-dialog"

  orderPositions: Signal<OrderPosition[]> = signal([])
  sum: Signal<number> = signal(0.0)
  checkoutCard: WritableSignal<boolean> = signal(false)

  constructor(
    private shoppingCart: ShoppingCartService,
  ){}

  ngOnInit(): void {
    this.orderPositions = this.shoppingCart.getOrderItems
    this.sum = this.shoppingCart.getSum
  }

  removeArticle(article: Article) {
    this.shoppingCart.removeArticle(article)
  }

  addArticle(article: Article) {
    this.shoppingCart.addArticle(article)
  }

  showCheckout(withCard: boolean){
    this.checkoutCard.set(withCard)
    const dialog = document.getElementById(this.CHECKOUT_DIALOG_ID) as HTMLDialogElement
    dialog.showModal()
  }
}
