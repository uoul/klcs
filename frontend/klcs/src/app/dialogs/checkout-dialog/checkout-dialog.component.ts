import { Component, input, OnInit, output, Signal, signal } from '@angular/core';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { CommonModule } from '@angular/common';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { FormsModule } from '@angular/forms';
import { OrderPosition } from '../../domain/OrderPosition';

@Component({
  selector: 'klcs-checkout-dialog',
  imports: [
    CommonModule,
    FormsModule,
  ],
  templateUrl: './checkout-dialog.component.html',
  styleUrl: './checkout-dialog.component.css'
})
export class CheckoutDialogComponent {
  dialogId = input.required<string>();
  withCard = input<boolean>(false);

  dialogClosed = output<void>();

  accountId: string = "";
  description: string ="";

  constructor(
    protected shoppingCart: ShoppingCartService,
    private sellerApi: SellerApiService,
  ){}

  checkout() {
    if(this.withCard()){
      const sub = this.sellerApi.checkoutCard(this.accountId, this.description).subscribe({
        next: val => {
          console.log(JSON.stringify(val))
          this.shoppingCart.clearCart()
        },
        error: err => console.error(err),
        complete: () => sub.unsubscribe()
      })
    }
    else {
      const sub = this.sellerApi.checkoutCash(this.description).subscribe({
        next: val => {
          console.log(JSON.stringify(val))
          this.shoppingCart.clearCart()
        },
        error: err => console.error(err),
        complete: () => sub.unsubscribe()
      })
    }
  }
}
