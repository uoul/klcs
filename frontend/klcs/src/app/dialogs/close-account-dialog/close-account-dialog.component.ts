import { CommonModule } from '@angular/common';
import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {ZXingScannerModule} from "@zxing/ngx-scanner";
import { AccountDetails } from '../../domain/AccountDetails';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-close-account-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ZXingScannerModule,
    TranslatePipe,
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
  isClosing: WritableSignal<boolean> = signal(false)

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
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

  async closeAccount(){
    if(!this.isClosing()){
      this.isClosing.set(true)
      try {
        this.accountDetails.set(
          await firstValueFrom(this.accountManagerApi.closeAccount(this.accountId()))
        )
      } finally { this.isClosing.set(false) }
    }
  }
}
