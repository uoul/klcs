import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { AccountDetails } from '../../domain/AccountDetails';
import { CommonModule } from '@angular/common';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

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

  newAccountDetails: WritableSignal<AccountDetails|null> = signal(null)

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
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
    const sub = this.accountManagerApi.postToAccount(this.accountId(), this.amount() * 100).subscribe({
      next: val => this.newAccountDetails.set(val),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }
}
