import {Injectable, Signal, signal, WritableSignal} from '@angular/core';
import { Article } from '../../domain/Article';
import { OrderPosition } from '../../domain/OrderPosition';

@Injectable({
  providedIn: 'root'
})
export class ShoppingCartService {

  private orderItems: OrderPosition[] = [];

  private _orderItems: WritableSignal<OrderPosition[]> = signal(this.orderItems);
  private _locked: WritableSignal<boolean> = signal(false);
  private _isEmpty: WritableSignal<boolean> = signal(true);
  private _sum: WritableSignal<number> = signal(0);

  constructor() {}

  public addArticle(article: Article) {
    if (!this._locked()){
      const orderPosition = this.orderItems.find(pos => {
        if (article.Id === pos.article.Id){
          return pos;
        }
        return undefined;
      });
      if (orderPosition) {
        orderPosition.count++;
      } else {
        this.orderItems.push({
          article: article,
          count: 1,
        });
      }
      this.updateCartStatus();
    }
  }

  public removeArticle(article: Article) {
    if(!this._locked()){
      const orderPosition = this.orderItems.find(pos => {
        if (article.Id === pos.article.Id){
          return pos;
        }
        return undefined;
      });
      if (orderPosition) {
        orderPosition.count--;
        if (orderPosition.count <= 0) {
          const index = this.orderItems.indexOf(orderPosition);
          if (index >= 0)
            this.orderItems.splice(index, 1);
        }
      }
      this.updateCartStatus();
    }
  }

  public get isLocked(): Signal<boolean> {
    return this._locked;
  }

  public lock() {
    this._locked.set(true);
  }

  public unlock() {
    this._locked.set(false);
  }

  public get isEmpty(): Signal<boolean> {
    return this._isEmpty;
  }

  public get getSum(): Signal<number> {
    return this._sum;
  }

  public get getOrderItems(): Signal<OrderPosition[]> {
    return this._orderItems;
  }

  public getPreparedOrder(): {[name: string]: number} {
    const retVal: {[name: string]: number} = {};
    for (const item of this.orderItems) {
      retVal[item.article.Id] = item.count;
    }
    return retVal;
  }

  public clearCart() {
    if(!this._locked()){
      this.orderItems.splice(0, this.orderItems.length);
      this.updateCartStatus();
    }
  }

  private updateCartStatus() {
    this._isEmpty.set(this.orderItems.length <= 0);
    this._sum.set(this.calcSum());
  }

  private calcSum(): number {
    let sum: number = 0;
    for (const pos of this.orderItems) {
      sum += pos.count * pos.article.Price;
    }
    return sum;
  }
}
