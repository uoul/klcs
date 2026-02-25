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
import {
  catchError,
  firstValueFrom,
  from,
  map,
  NEVER,
  of,
  switchMap,
  tap,
} from 'rxjs';
import { PublicApiService } from './services/public-api/public-api.service';
import { AuthService } from './services/auth/auth.service';

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
      const authService = inject(AuthService);

      return firstValueFrom(
        publicApiService.getSettings().pipe(
          tap((appSettings) => {
            oauthService.configure({
              issuer: appSettings.Oidc.Authority,
              clientId: appSettings.Oidc.ClientId,
              redirectUri: window.location.origin,
              scope: 'openid profile email offline_access',
              responseType: 'code',
              sessionChecksEnabled: false,
            });
            oauthService.setStorage(localStorage);
          }),
          switchMap(() =>
            from(oauthService.loadDiscoveryDocumentAndTryLogin()),
          ),
          switchMap(() => {
            if (oauthService.hasValidAccessToken()) {
              oauthService.setupAutomaticSilentRefresh();
              return of(true);
            }

            if (oauthService.getRefreshToken()) {
              return from(oauthService.refreshToken()).pipe(
                tap(() => oauthService.setupAutomaticSilentRefresh()),
                map(() => true),
                catchError(() => {
                  oauthService.initLoginFlow();
                  return NEVER;
                }),
              );
            }

            oauthService.initLoginFlow();
            return NEVER;
          }),
          tap(() => authService.setReady()),
        ),
      );
    }),
  ],
};
