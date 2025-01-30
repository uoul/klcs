import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, of } from 'rxjs';
import { Shop } from '../../domain/Shop';
import { KlcsConfig } from '../../config/KlcsConfig';

@Injectable({
  providedIn: 'root'
})
export class KlcsAdminApiService {

  constructor(
    private http: HttpClient,
  ) { }

  public getShops(): Observable<Shop[]> {
    return this.http.get<Shop[]>(`${KlcsConfig.BackendRoot}/api/v1/admin/shops`);
  }

  public createShop(shop: Shop): Observable<Shop>{
    return this.http.post<Shop>(`${KlcsConfig.BackendRoot}/api/v1/admin/shops`, shop);
  }

  public updateShop(shop: Shop): Observable<Shop>{
    return this.http.patch<Shop>(`${KlcsConfig.BackendRoot}/api/v1/admin/shops/${shop.Id}`, shop);
  }

  public deleteShop(shopId: string): Observable<Object> {
    return this.http.delete(`${KlcsConfig.BackendRoot}/api/v1/admin/shops/${shopId}`);
  }
}
