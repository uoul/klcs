import { inject } from '@angular/core';
import { CanActivateFn } from '@angular/router';
import { OAuthService } from 'angular-oauth2-oidc';

export const authGuard: CanActivateFn = () => {
  const authService: OAuthService = inject(OAuthService);
  return authService.loadDiscoveryDocumentAndTryLogin()
    .then((_) => {
      return authService.hasValidAccessToken() && authService.hasValidIdToken();
    });
};
