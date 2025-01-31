import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {catchError, finalize, Observable, of, tap, throwError} from "rxjs";
import {ShoppingCartService} from "../shopping-cart/shopping-cart.service";
import { KlcsConfig } from '../../config/KlcsConfig';
import { Order } from '../../domain/Order';
import { Shop } from '../../domain/Shop';
import { ShopDetails } from '../../domain/ShopDetails';

@Injectable({
  providedIn: 'root'
})
export class SellerApiService {

  constructor(
    private http: HttpClient,
    private cartService: ShoppingCartService,
  ) { }

  public getShops(): Observable<Shop[]> {
    return this.http.get<Shop[]>(`${KlcsConfig.BackendRoot}/api/v1/shops`);
  }

  public getShopDetails(shopId: string): Observable<ShopDetails> {
    return this.http.get<ShopDetails>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}`);
  }

  public checkoutCard(accountId: string, description: string): Observable<Order> {
    const order: Order = {
      AccountId: accountId,
      Type: "CARD",
      Description: description,
      Articles: this.cartService.getPreparedOrder(),
      Sum: undefined,
    }
    return this.placeOrder(order);
  }

  public checkoutCash(description: string): Observable<Order> {
    const items = this.cartService.getPreparedOrder();
    const order: Order = {
      AccountId: null,
      Type: "CASH",
      Description: description,
      Articles: items,
      Sum: undefined,
    }
    return this.placeOrder(order);
  }

  private placeOrder(order: Order): Observable<Order> {
    this.cartService.lock();
    return this.http.post<Order>(`${KlcsConfig.BackendRoot}/api/v1/orders`, order).pipe(
      finalize(() => {
        this.cartService.unlock();
      }),
    );
  }
}
