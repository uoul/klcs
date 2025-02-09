import { Component, input, output } from '@angular/core';
import { ShoppingCartComponent } from "../../components/shopping-cart/shopping-cart.component";
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

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
}
