import { Component, computed, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { Account } from '../../domain/Account';
import { FormsModule } from '@angular/forms';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { TranslatePipe, TranslateService, TranslateDirective } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-create-account-dialog',
  imports: [
    FormsModule,
    TranslatePipe,
],
  templateUrl: './create-account-dialog.component.html',
  styleUrl: './create-account-dialog.component.css'
})
export class CreateAccountDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output();
  accountCreated: OutputEmitterRef<Account> = output();

  _account: WritableSignal<Account> = signal(new Account())
  _lockDataTip = computed(()=> this._account().Locked ? "Unlock" : "Lock")
  isActive: WritableSignal<boolean> = signal(false)

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  _dialogClosed() {
    this.init()
    this.dialogClosed.emit()
  }

  init() {
    this._account.set(new Account())
  }

  async createAccount(){
    if(!this.isActive()){
      this.isActive.set(true)
      try {
        const created = await firstValueFrom(this.accountManagerApi.createAccount(this._account()))
        this.accountCreated.emit(created)
      } finally { this.isActive.set(false) }
    }
  }
}
