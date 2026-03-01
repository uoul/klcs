import { CommonModule } from '@angular/common';
import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AccountDetails } from '../../domain/AccountDetails';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-check-account-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule,
    TranslatePipe,
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
  isCheckActive: WritableSignal<boolean> = signal(false)

  json = JSON;

  constructor(
    private accountManagerApi: AccountManagerApiService,
    protected translate: TranslateService,
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

  async checkAccount() {
    if(!this.isCheckActive()){
      this.isCheckActive.set(true)
      try {
        const details = await firstValueFrom(this.accountManagerApi.getAccountDetails(this.accountId))
        this.accountData.set(details)
      } finally { this.isCheckActive.set(false) }
    }
  }
}
