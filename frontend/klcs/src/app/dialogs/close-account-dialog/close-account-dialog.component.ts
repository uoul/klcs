import { CommonModule } from '@angular/common';
import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { AccountDetails } from '../../domain/AccountDetails';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-close-account-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule
  ],
  templateUrl: './close-account-dialog.component.html',
  styleUrl: './close-account-dialog.component.css'
})
export class CloseAccountDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output();
    
  accountId: WritableSignal<string> = signal("")
  scannerActive: WritableSignal<boolean> = signal(false)
  accountDetails: WritableSignal<AccountDetails|null> = signal(null)

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
  ){}

  showScanner(){
    this.scannerActive.set(true)
  }

  _dialogClosed(){
    this.scannerActive.set(false)
    this.accountDetails.set(null)
    this.accountId.set("")
  }

  onScanSuccess(data: string) {
    this.accountId.set(data)
    this.scannerActive.set(false)
    this.closeAccount()
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive.set(false)
  }

  closeAccount(){
    const sub = this.accountManagerApi.closeAccount(this.accountId()).subscribe({
      next: val => this.accountDetails.set(val),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }
}
