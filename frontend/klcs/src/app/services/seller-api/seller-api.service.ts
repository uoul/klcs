import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Shop } from '../../domain/Shop';
import { KlcsConfig } from '../../config/KlcsConfig';

@Injectable({
  providedIn: 'root'
})
export class SellerApiService {

  constructor(
    private http: HttpClient,
  ) { }

  public getShops(): Observable<Shop[]> {
    return this.http.get<Shop[]>(`${KlcsConfig.BackendRoot}/api/v1/shops`)
  }
}
