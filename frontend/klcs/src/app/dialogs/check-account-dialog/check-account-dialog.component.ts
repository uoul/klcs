import { CommonModule } from '@angular/common';
import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AccountDetails } from '../../domain/AccountDetails';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import {ZXingScannerModule} from "@zxing/ngx-scanner";

@Component({
  selector: 'klcs-check-account-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule,
  ],
  templateUrl: './check-account-dialog.component.html',
  styleUrl: './check-account-dialog.component.css'
})
export class CheckAccountDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output();
  
  scannerActive: boolean = false
  accountId: string = ""
  accountData: WritableSignal<AccountDetails|null> = signal(null)

  json = JSON;

  constructor(
    private accountManagerApi: AccountManagerApiService,
  ){}

  _dialogClosed(){
    this.accountId = "";
    this.scannerActive = false;
    this.accountData.set(null);
    this.dialogClosed.emit();
  }

  onScanSuccess(data: string) {
    this.accountId = data;
    this.scannerActive = false;
    this.checkAccount();
  }

  onScanError(error: Error){
    console.error(error)
    this.scannerActive = false;
  }

  checkAccount() {
    const sub = this.accountManagerApi.getAccountDetails(this.accountId).subscribe({
      next: details => this.accountData.set(details),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
