import { Component, computed, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { Account } from '../../domain/Account';
import { FormsModule } from '@angular/forms';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-create-account-dialog',
  imports: [
    FormsModule,
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

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
  ){}

  _dialogClosed() {
    this.init()
    this.dialogClosed.emit()
  }

  init() {
    this._account.set(new Account())
  }

  createAccount(){
    const sub = this.accountManagerApi.createAccount(this._account()).subscribe({
      next: val => this.accountCreated.emit(val),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }
}
