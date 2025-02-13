import { Component, input, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Account } from '../../domain/Account';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';

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
  ){}

  _dialogClosed() {
    this.dialogClosed.emit()
  }

  updateAccount(){
    const sub = this.accountManagerApi.updateAccount(this.account()).subscribe({
      next: _ => this.accountUpdated.emit(this.account()),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
