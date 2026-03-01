import { Component, input, output, signal, WritableSignal } from '@angular/core';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { CommonModule } from '@angular/common';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { FormsModule } from '@angular/forms';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { FocusDirective } from '../../directives/focus/focus.directive';
import {  TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-checkout-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule,
    FocusDirective,
    TranslatePipe,
  ],
  templateUrl: './checkout-dialog.component.html',
  styleUrl: './checkout-dialog.component.css'
})
export class CheckoutDialogComponent {
  dialogId = input.required<string>();
  withCard = input<boolean>(false);

  dialogClosed = output<void>();

  accountId: WritableSignal<string> = signal("");
  description: WritableSignal<string> = signal("");
  scannerActive: WritableSignal<boolean> = signal(false);
  isActiveCheckout: WritableSignal<boolean> = signal(false)

  constructor(
    protected shoppingCart: ShoppingCartService,
    private sellerApi: SellerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  _dialogClosed(){
    this.accountId.set("")
    this.description.set("")
    this.scannerActive.set(false)
    this.dialogClosed.emit();
  }

  async checkout() {
    if(!this.isActiveCheckout()) {
      this.isActiveCheckout.set(true)
      if(this.withCard()){
        try {
          await firstValueFrom(this.sellerApi.checkoutCard(this.accountId(), this.description()))
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: this.translate.instant("success.OrderPlaced")})
          await this.sellerApi.refreshShopDetails()
        } finally { this.isActiveCheckout.set(false) }
      }
      else {
        try {
          await firstValueFrom(this.sellerApi.checkoutCash(this.description()))
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: this.translate.instant("success.OrderPlaced")})
          await  this.sellerApi.refreshShopDetails()
        } finally { this.isActiveCheckout.set(false) }
      }
    }
  }

  onScanSuccess(data: string) {
    this.accountId.set(data);
    this.scannerActive.set(false);
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive.set(false);
  }
}
