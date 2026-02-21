import { Component, input, OnInit, output, signal, WritableSignal } from '@angular/core';
import { ShoppingCartComponent } from "../../components/shopping-cart/shopping-cart.component";
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { FocusDirective } from '../../directives/focus/focus.directive';
import { Article } from '../../domain/Article';
import { PaymentItemListComponent } from "../../components/payment-item-list/payment-item-list.component";
import { PublicApiService } from '../../services/public-api/public-api.service';

@Component({
  selector: 'klcs-edit-cart-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ShoppingCartComponent,
    ZXingScannerModule,
    FocusDirective,
    PaymentItemListComponent
],
  templateUrl: './edit-cart-dialog.component.html',
  styleUrl: './edit-cart-dialog.component.css'
})
export class EditCartDialogComponent implements OnInit {
  dialogId = input.required<string>()
  dialogClosed = output<void>()

  paymentMethod: WritableSignal<number> = signal(1);
  accountId: WritableSignal<string> = signal("");
  description: WritableSignal<string> = signal("");
  scannerActive: WritableSignal<"none" | "accountId" | "description"> = signal("none");
  step: WritableSignal<"order" | "payment"> = signal("order")

  public paymentItems: WritableSignal<Article[]> = signal([])

  constructor(
    protected publicApi: PublicApiService,
    protected shoppingCart: ShoppingCartService,
    private sellerApi: SellerApiService,
    private notify: NotificationService,
  ){}

  ngOnInit(): void {
    // Set Default Payment Method
    if(this.publicApi.settings()?.UiSettings.Mobile.DefaultPayment == "CARD") {
      this.paymentMethod.set(2)
    }
  }

  checkout() {
    if(this.paymentMethod() == 2){
      const sub = this.sellerApi.checkoutCard(this.accountId(), this.description()).subscribe({
        next: _ => {
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: "Successfully placed order"})
          this.sellerApi.refreshShopDetails()
          this.close()
        },
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
        complete: () => sub.unsubscribe()
      })
    }
    else if(this.paymentMethod() == 1) {
      this.paymentItems.set(this.shoppingCart.getItems())
      const sub = this.sellerApi.checkoutCash(this.description()).subscribe({
        next: _ => {
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: "Successfully placed order"})
          this.sellerApi.refreshShopDetails()
          this.step.set("payment")
        },
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
        complete: () => sub.unsubscribe()
      })
    }
  }

  _dialogClosed(){
    this.paymentMethod.set(1);
    this.accountId.set("");
    this.description.set("");
    this.scannerActive.set("none");
    this.step.set("order");
    this.paymentItems.set([]);
    this.dialogClosed.emit();
  }

  onScanSuccess(data: string) {
    switch(this.scannerActive()) {
      case "accountId":
        this.accountId.set(data);
        break;
      case "description":
        this.description.set(data)
        break;
    }
    this.scannerActive.set("none");
  }

  close() {
    const dialog = document.getElementById(this.dialogId()) as HTMLDialogElement
    dialog.close()
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive.set("none");
  }

}
