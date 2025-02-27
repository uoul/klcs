import { Component, input, output } from '@angular/core';
import { ShoppingCartComponent } from "../../components/shopping-cart/shopping-cart.component";
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-edit-cart-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ShoppingCartComponent,
    ZXingScannerModule,
  ],
  templateUrl: './edit-cart-dialog.component.html',
  styleUrl: './edit-cart-dialog.component.css'
})
export class EditCartDialogComponent {
  dialogId = input.required<string>()
  dialogClosed = output<void>()

  paymentMethod: number = 0;
  accountId: string = "";
  description: string = "";
  scannerActive: boolean = false;

  constructor(
    protected shoppingCart: ShoppingCartService,
    private sellerApi: SellerApiService,
    private notify: NotificationService,
  ){}

  checkout() {
    if(this.paymentMethod == 2){
      const sub = this.sellerApi.checkoutCard(this.accountId, this.description).subscribe({
        next: _ => this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: "Successfully placed order"}),
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
        complete: () => sub.unsubscribe()
      })
    }
    else if(this.paymentMethod == 1) {
      const sub = this.sellerApi.checkoutCash(this.description).subscribe({
        next: _ => this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: "Successfully placed order"}),
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
        complete: () => sub.unsubscribe()
      })
    }
  }

  _dialogClosed(){
    this.paymentMethod = 0;
    this.accountId = "";
    this.description = "";
    this.scannerActive = false;
    this.dialogClosed.emit();
  }

  onScanSuccess(data: string) {
    this.accountId = data;
    this.scannerActive = false;
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive = false;
  }

}
