import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { ZXingScannerModule } from "@zxing/ngx-scanner";
import { KlcsConfig } from '../../config/KlcsConfig';
import { AccountDetails } from '../../domain/AccountDetails';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';

@Component({
  selector: 'klcs-charge-account-dialog',
  imports: [
    CommonModule,
    ZXingScannerModule,
    FormsModule,
  ],
  templateUrl: './charge-account-dialog.component.html',
  styleUrl: './charge-account-dialog.component.css'
})
export class ChargeAccountDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output();
  
  scannerActive: WritableSignal<boolean> = signal(false)
  accountId: WritableSignal<string> = signal("")
  amount: WritableSignal<number> = signal(0)
  chargeActive: WritableSignal<boolean> = signal(false)

  newAccountDetails: WritableSignal<AccountDetails|null> = signal(null)

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  onScanSuccess(data: string) {
    this.accountId.set(data)
    this.scannerActive.set(false)
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive.set(false)
  }

  _dialogClosed() {
    this.dialogClosed.emit()
    this.scannerActive.set(false)
    this.accountId.set("")
    this.amount.set(0)
    this.newAccountDetails.set(null)
  }

  showScanner(){
    this.scannerActive.set(true)
  }

  chargeAccount(){
    if(!this.chargeActive()){
      this.chargeActive.set(true)
      const sub = this.accountManagerApi.postToAccount(this.accountId(), this.amount() * 100).subscribe({
        next: val => this.newAccountDetails.set(val),
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => {
          this.chargeActive.set(false)
          sub.unsubscribe()
        },
      })
    } else {
      this.notify.show({type: "warning", duration: KlcsConfig.durationError, message: this.translate.instant("warnings.WarnChargeTwice")})
    }
  }
}
