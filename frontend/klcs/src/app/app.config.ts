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
import { firstValueFrom, switchMap, tap } from 'rxjs';
import { PublicApiService } from './services/public-api/public-api.service';

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
          // Changed from getShops()
          tap((appSettings) => {
            oauthService.configure({
              issuer: appSettings.Oidc.Authority,
              clientId: appSettings.Oidc.ClientId,
              redirectUri: window.location.origin,
              scope: 'openid profile email',
              responseType: 'code',
            });
            oauthService.setStorage(localStorage);
            oauthService.setupAutomaticSilentRefresh();
          }),
          switchMap(() => oauthService.loadDiscoveryDocumentAndLogin()),
        ),
      );
    }),
  ],
};
