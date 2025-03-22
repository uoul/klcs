import {Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {finalize, Observable, subscribeOn, tap} from "rxjs";
import {ShoppingCartService} from "../shopping-cart/shopping-cart.service";
import { KlcsConfig } from '../../config/KlcsConfig';
import { Order } from '../../domain/Order';
import { Shop } from '../../domain/Shop';
import { ShopDetails } from '../../domain/ShopDetails';
import { NotificationService } from '../notification/notification.service';
import { HistoryItem } from '../../domain/HistoryItem';
import { Article } from '../../domain/Article';

@Injectable({
  providedIn: 'root'
})
export class SellerApiService {

  constructor(
    private http: HttpClient,
    private cartService: ShoppingCartService,
    private notify: NotificationService,
  ) { }

  _shopId: WritableSignal<string> = signal("")
  _shopDetails: WritableSignal<ShopDetails> = signal(new ShopDetails())

  public getShops(): Observable<Shop[]> {
    return this.http.get<Shop[]>(`${KlcsConfig.BackendRoot}/api/v1/shops`);
  }

  public get getShopDetails(): Signal<ShopDetails> {
    return this._shopDetails
  }

  public updateShopId(shopId: string): void {
    this._shopId.set(shopId)
    this.refreshShopDetails()
  }

  public refreshShopDetails(): void {
    if(this._shopId().length > 0) {
      const sub = this.http.get<ShopDetails>(`${KlcsConfig.BackendRoot}/api/v1/shops/${this._shopId()}`).subscribe({
        next: s => {
          for(let [name, articles] of Object.entries(s.Categories)){
            articles.sort((a: Article, b: Article) => a.Name == b.Name  ? 0 : a.Name < b.Name ? -1 : 1 )
          }
          this._shopDetails.set(s)
        },
        error: err => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
        complete: () => sub.unsubscribe()
      })
    }
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

  public getHistory(len: number): Observable<HistoryItem[]> {
    return this.http.get<HistoryItem[]>(`${KlcsConfig.BackendRoot}/api/v1/history?length=${len}`).pipe(
      tap(history => {
        for(let h of history){
          h.Timestamp = new Date(h.Timestamp)
        }
      })
    )
  }

  private placeOrder(order: Order): Observable<Order> {
    this.cartService.lock();
    return this.http.post<Order>(`${KlcsConfig.BackendRoot}/api/v1/orders`, order).pipe(
      finalize(() => {
        this.cartService.unlock();
        this.cartService.clearCart();
      }),
    );
  }
}
