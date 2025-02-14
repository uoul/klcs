import { Component, input, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Account } from '../../domain/Account';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-edit-account-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './edit-account-dialog.component.html',
  styleUrl: './edit-account-dialog.component.css'
})
export class EditAccountDialogComponent {
  dialogId = input.required<string>()
  dialogClosed = output<void>()
  accountUpdated = output<Account>()

  account = input.required<Account>()

  constructor(
    private accountManagerApi: AccountManagerApiService,
    private notify: NotificationService,
  ){}

  _dialogClosed() {
    this.dialogClosed.emit()
  }

  updateAccount(){
    const sub = this.accountManagerApi.updateAccount(this.account()).subscribe({
      next: _ => this.accountUpdated.emit(this.account()),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }
}
