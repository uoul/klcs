import { HttpInterceptorFn } from '@angular/common/http';
import { AuthService } from '../../services/auth/auth.service';
import { inject } from '@angular/core';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService: AuthService = inject(AuthService);
  if (!authService.isLoggedIn()) {
    authService.login();
  }
  return next(req.clone({
    setHeaders: {
      Authorization: `Bearer ${authService.getAccessToken()}`
    }
  }));
};
