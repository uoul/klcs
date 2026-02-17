import { Component, input, output } from '@angular/core';
import { OrderPosition } from '../../domain/OrderPosition';

@Component({
  selector: 'klcs-shopping-cart-item',
  imports: [],
  templateUrl: './shopping-cart-item.component.html',
  styleUrl: './shopping-cart-item.component.css'
})
export class ShoppingCartItemComponent {
  orderPosition = input.required<OrderPosition>()
  decrementClicked = output<OrderPosition>()
  incrementClicked = output<OrderPosition>()
}
