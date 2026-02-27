import { Component, input, output, signal, WritableSignal } from '@angular/core';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { CommonModule } from '@angular/common';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { FormsModule } from '@angular/forms';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { ErrorResponse } from '../../domain/ErrorResponse';
import { FocusDirective } from '../../directives/focus/focus.directive';
import { TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'klcs-checkout-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule,
    FocusDirective,
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

  checkout() {
    if(this.withCard()){
      const sub = this.sellerApi.checkoutCard(this.accountId(), this.description()).subscribe({
        next: _ => {
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: this.translate.instant("success.OrderPlaced")})
          this.sellerApi.refreshShopDetails()
        },
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => sub.unsubscribe()
      })
    }
    else {
      const sub = this.sellerApi.checkoutCash(this.description()).subscribe({
        next: _ => {
          this.notify.show({type: "success", duration: KlcsConfig.durationShort, message: this.translate.instant("success.OrderPlaced")})
          this.sellerApi.refreshShopDetails()
        },
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => sub.unsubscribe()
      })
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
