import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { catchError, throwError } from 'rxjs';
import { KlcsConfig } from '../../config/KlcsConfig';
import { NotificationService } from '../../services/notification/notification.service';
import { AuthService } from '../../services/auth/auth.service';

export const notifyInterceptor: HttpInterceptorFn = (req, next) => {
  const notify = inject(NotificationService);
  const translate = inject(TranslateService);
  const auth = inject(AuthService)

  return next(req).pipe(
    catchError((err: HttpErrorResponse) => {
      // Handle globally — e.g. 401, 403, 500
      if (err.status === 401) {
        auth.login()
      } else {
        notify.show({
          type: 'error',
          duration: KlcsConfig.durationError,
          message: translate.instant(`errors.${err.error?.Code}`),
        });
      }
      return throwError(() => err);
    })
  );
};
