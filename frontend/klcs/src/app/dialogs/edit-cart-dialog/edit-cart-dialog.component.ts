import { Component, input, output } from '@angular/core';
import { ShoppingCartComponent } from "../../components/shopping-cart/shopping-cart.component";
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../services/shopping-cart/shopping-cart.service';

@Component({
  selector: 'klcs-edit-cart-dialog',
  imports: [
    CommonModule,
    FormsModule,
    ShoppingCartComponent,
  ],
  templateUrl: './edit-cart-dialog.component.html',
  styleUrl: './edit-cart-dialog.component.css'
})
export class EditCartDialogComponent {
  dialogId = input.required<string>()
  dialogClosed = output<void>()

  paymentMethod: number = 0;
  accountId: string = "";

  constructor(
    protected shoppingCart: ShoppingCartService,
  ){}

  checkOut() {
    console.log("checkout")
  }

  _dialogClosed(){
    this.paymentMethod = 0;
    this.accountId = "";
    this.dialogClosed.emit();
  }
}
