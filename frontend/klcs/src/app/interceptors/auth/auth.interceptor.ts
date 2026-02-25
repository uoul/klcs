import { HttpInterceptorFn } from '@angular/common/http';
import { AuthService } from '../../services/auth/auth.service';
import { inject } from '@angular/core';
import { switchMap } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService: AuthService = inject(AuthService);
  // Skip waiting for Keycloak/settings requests to avoid deadlock
  if (req.url.includes('/realms/') || req.url.includes('/settings')) {
    return next(req);
  }

  return authService.waitForReady().pipe(
    switchMap(() =>
      next(
        req.clone({
          setHeaders: {
            Authorization: `Bearer ${authService.getAccessToken()}`,
          },
        }),
      ),
    ),
  );
};
