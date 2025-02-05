import { Component, OnInit, signal, Signal } from '@angular/core';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';
import { OrderPosition } from '../../domain/OrderPosition';
import { ShoppingCartItemComponent } from "../shopping-cart-item/shopping-cart-item.component";
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'klcs-shopping-cart',
  imports: [
    CommonModule,
    ShoppingCartItemComponent,
  ],
  templateUrl: './shopping-cart.component.html',
  styleUrl: './shopping-cart.component.css'
})
export class ShoppingCartComponent implements OnInit {
  
  orderPositions: Signal<OrderPosition[]> = signal([])
  sum: Signal<number> = signal(0.0)

  constructor(
    private shoppingCart: ShoppingCartService,
  ){}

  ngOnInit(): void {
    this.orderPositions = this.shoppingCart.getOrderItems
    this.sum = this.shoppingCart.getSum
  }

  removeArticle(article: Article) {
    this.shoppingCart.removeArticle(article)
  }

  addArticle(article: Article) {
    this.shoppingCart.addArticle(article)
  }
}
