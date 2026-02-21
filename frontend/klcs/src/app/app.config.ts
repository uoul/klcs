import {
  ApplicationConfig,
  provideZoneChangeDetection,
  isDevMode,
  provideAppInitializer,
  inject,
} from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideServiceWorker } from '@angular/service-worker';
import { OAuthService, provideOAuthClient } from 'angular-oauth2-oidc';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor } from './interceptors/auth/auth.interceptor';
import { provideAnimations } from '@angular/platform-browser/animations';
import { catchError, firstValueFrom, from, of, switchMap, tap } from 'rxjs';
import { PublicApiService } from './services/public-api/public-api.service';
import { animate } from '@angular/animations';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideServiceWorker('ngsw-worker.js', {
      enabled: !isDevMode(),
      registrationStrategy: 'registerWhenStable:30000',
    }),
    provideOAuthClient(),
    provideHttpClient(withInterceptors([authInterceptor])),
    provideAnimations(),
    provideAppInitializer(() => {
      const publicApiService = inject(PublicApiService);
      const oauthService = inject(OAuthService);
      return firstValueFrom(
        publicApiService.getSettings().pipe(
          tap((appSettings) => {
            oauthService.configure({
              issuer: appSettings.Oidc.Authority,
              clientId: appSettings.Oidc.ClientId,
              redirectUri: window.location.origin,
              scope: 'openid profile email',
              responseType: 'code',
              sessionChecksEnabled: false, // avoids iframe issues with Keycloak
            });
            oauthService.setStorage(localStorage);
          }),
          switchMap(() => oauthService.loadDiscoveryDocumentAndTryLogin()),
          switchMap(() => {
            if (oauthService.hasValidAccessToken()) {
              // Token still valid — start refresh timer
              oauthService.setupAutomaticSilentRefresh();
              return of(true);
            }
            if (oauthService.getRefreshToken()) {
              // Access token expired but refresh token exists — try to recover silently
              return from(oauthService.refreshToken()).pipe(
                tap(() => oauthService.setupAutomaticSilentRefresh()),
                catchError(() => {
                  // Refresh token also expired — force re-login
                  oauthService.initLoginFlow();
                  return of(false);
                }),
              );
            }
            // No tokens at all — redirect to login
            oauthService.initLoginFlow();
            return of(false);
          }),
        ),
      );
    }),
  ],
};
