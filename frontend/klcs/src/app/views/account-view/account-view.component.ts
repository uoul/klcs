import { Component, computed, OnInit, signal, WritableSignal } from '@angular/core';
import { AccountManagerApiService } from '../../services/account-manager-api/account-manager-api.service';
import { Account } from '../../domain/Account';
import { CreateAccountDialogComponent } from "../../dialogs/create-account-dialog/create-account-dialog.component";
import { FormsModule } from '@angular/forms';
import { ReadQrDialogComponent } from "../../dialogs/read-qr-dialog/read-qr-dialog.component";
import { EditAccountDialogComponent } from "../../dialogs/edit-account-dialog/edit-account-dialog.component";

@Component({
  selector: 'klcs-account-view',
  imports: [
    FormsModule,
    CreateAccountDialogComponent,
    ReadQrDialogComponent,
    EditAccountDialogComponent
],
  templateUrl: './account-view.component.html',
  styleUrl: './account-view.component.css'
})
export class AccountViewComponent implements OnInit {

  accounts: WritableSignal<Account[] | null> = signal(null)
  searchText: WritableSignal<string> = signal("")
  selectedAccount: WritableSignal<Account> = signal(new Account())

  filteredAccounts = computed(() => {
    return this.accounts()?.filter(a => {
      try {
        return `${a.Id}${a.HolderName}`.toLowerCase().match(this.searchText().toLowerCase()) === null ? false : true
      } catch (e) {
        return true
      }
    }).sort((a: Account, b: Account) => {
      if(a == b)
        return 0;
      if(a.HolderName < b.HolderName)
        return -1
      return 1
    })
  })

  protected readonly CREATE_ACCOUNT_DIALOG_ID = "create-account-dialog"
  protected readonly READ_QR_DIALOG_ID = "read-qr-for-account-dialog"
  protected readonly EDIT_ACCOUNT_DIALOG_ID = "edit-account-dialog"

  constructor(
    private accountManagerApi: AccountManagerApiService,
  ){}

  ngOnInit(): void {
    this.refresh()
  }

  refresh() {
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
        next: _ => this.refresh(),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      })
    }
  }

  showCreateAccountDialog(){
    const dialog = document.getElementById(this.CREATE_ACCOUNT_DIALOG_ID) as HTMLDialogElement
    dialog.showModal()
  }

  showReadQrDialog(){
    const dialog = document.getElementById(this.READ_QR_DIALOG_ID) as HTMLDialogElement
    dialog.showModal() 
  }

  showEditAccountDialog(account: Account){
    const dialog = document.getElementById(this.EDIT_ACCOUNT_DIALOG_ID) as HTMLDialogElement
    this.selectedAccount.set(JSON.parse(JSON.stringify(account)))
    dialog.showModal() 
  }

  filter(){}
}
