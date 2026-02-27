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
import {provideTranslateService, provideTranslateLoader} from "@ngx-translate/core";
import {provideTranslateHttpLoader} from "@ngx-translate/http-loader";
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
    provideTranslateService({
      loader: provideTranslateHttpLoader({
        prefix: "/i18n/",
        suffix: ".json"
      }),
      fallbackLang: "en",
    }),
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
              scope: 'openid profile email',
              responseType: 'code',
              sessionChecksEnabled: false,
              clockSkewInSec: 60,
            });
            oauthService.setStorage(localStorage);
          }),
          switchMap(() =>
            from(oauthService.loadDiscoveryDocumentAndTryLogin()),
          ),
          switchMap(() => {
            const expiresAt = oauthService.getAccessTokenExpiration();
            const isActuallyValid = oauthService.hasValidAccessToken() && expiresAt > Date.now();

            if (isActuallyValid) {
              return of(true);
            }

            if (oauthService.getRefreshToken()) {
              return from(oauthService.refreshToken()).pipe(
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
          tap(() => {
            oauthService.setupAutomaticSilentRefresh();
            authService.setReady();                     
          }),
        ),
      );
    }),
  ],
};
