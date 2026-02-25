import { Injectable } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, Observable, filter, map, startWith, take } from 'rxjs';
import { UserIdentity } from '../../domain/UserIdentity';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(
    private oauthService: OAuthService
  ) { }

  private ready$ = new BehaviorSubject<boolean>(false);

  setReady(): void {
    this.ready$.next(true);
  }

  waitForReady() {
    return this.ready$.pipe(
      filter(ready => ready),
      take(1),
    );
  }

  login() {
    this.oauthService.initCodeFlow();
  }

  logout() {
    this.oauthService.revokeTokenAndLogout();
  }

  getAccessToken(): string {
    return this.oauthService.getAccessToken();
  }

  getIdentity(): UserIdentity {
    return UserIdentity.of(this.oauthService.getAccessToken());
  }

  isLoggedIn(): Observable<boolean> {
    return this.oauthService.events.pipe(
      map(_ => this.oauthService.hasValidAccessToken() && this.oauthService.hasValidIdToken()),
      startWith(this.oauthService.hasValidAccessToken() && this.oauthService.hasValidIdToken()),
    );
  }
}
