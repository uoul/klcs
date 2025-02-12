import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import { Account } from '../../domain/Account';
import { AccountDetails } from '../../domain/AccountDetails';
import { KlcsConfig } from '../../config/KlcsConfig';

@Injectable({
  providedIn: 'root'
})
export class AccountManagerApiService {

  constructor(
    private http: HttpClient,
  ) { }

  public getAccounts(): Observable<Account[]> {
    return this.http.get<Account[]>(`${KlcsConfig.BackendRoot}/api/v1/accounts`)
  }

  public getAccountDetails(accountId: string): Observable<AccountDetails> {
    return this.http.get<AccountDetails>(`${KlcsConfig.BackendRoot}/api/v1/accounts/${accountId}`)
  }

  public createAccount(account: Account): Observable<Account> {
    return this.http.post<Account>(`${KlcsConfig.BackendRoot}/api/v1/account`, account)
  }

  public updateAccount(account: Account): Observable<Account> {
    return this.http.patch<Account>(`${KlcsConfig.BackendRoot}/api/v1/accounts/${account.Id}`, account)
  }

  public closeAccount(accountId: string): Observable<AccountDetails> {
    return this.http.delete<AccountDetails>(`${KlcsConfig.BackendRoot}/api/v1/accounts/${accountId}/balance`)
  }

  public postToAccount(accountId: string, amount: number): Observable<AccountDetails> {
    return this.http.post<AccountDetails>(`${KlcsConfig.BackendRoot}/api/v1/accounts/${accountId}/balance`, { Amount: amount })
  }
}
