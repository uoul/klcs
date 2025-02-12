import { Component, OnInit, signal, WritableSignal } from '@angular/core';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { Account } from '../../domain/Account';

@Component({
  selector: 'klcs-account-view',
  imports: [],
  templateUrl: './account-view.component.html',
  styleUrl: './account-view.component.css'
})
export class AccountViewComponent implements OnInit {

  accounts: WritableSignal<Account[] | null> = signal(null)

  constructor(
    private accountManagerApi: AccountManagerApiService,
  ){}

  ngOnInit(): void {
    const sub = this.accountManagerApi.getAccounts().subscribe({
      next: a => this.accounts.set(a),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  setLocked(account: Account, state: boolean) {
    if(confirm(`Do you realy want to ${state ? "lock" : "unlock"} account ${account.Id}?`)){
      account.Locked = state
      const sub = this.accountManagerApi.updateAccount(account).subscribe({
        next: _ => console.log("Account has been updated"),
        error: err => {
          account.Locked = !account.Locked
          console.error(err)
        },
        complete: () => sub.unsubscribe(),
      })
    }
  }
}
