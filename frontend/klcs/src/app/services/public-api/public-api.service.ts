import { HttpClient } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { Observable, of, tap } from 'rxjs';
import { AppSettings } from '../../domain/AppSettings';
import { KlcsConfig } from '../../config/KlcsConfig';

@Injectable({
  providedIn: 'root',
})
export class PublicApiService {
  private appSettingsSignal = signal<AppSettings | undefined>(undefined);
  public readonly settings = this.appSettingsSignal.asReadonly();

  constructor(private http: HttpClient) { }

  public getSettings(): Observable<AppSettings> {
    return this.http.get<AppSettings>(`${KlcsConfig.BackendRoot}/public/settings`).pipe(
      tap(settings => this.appSettingsSignal.set(settings))
    );
  }
}
