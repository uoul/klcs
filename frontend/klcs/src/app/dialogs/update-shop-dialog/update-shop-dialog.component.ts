import { Component, EventEmitter, input, Input, InputSignal, OnInit, output, Output, OutputEmitterRef, Signal } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { FormsModule } from '@angular/forms';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';

@Component({
  selector: 'klcs-update-shop-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './update-shop-dialog.component.html',
  styleUrl: './update-shop-dialog.component.css'
})
export class UpdateShopDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  shop: InputSignal<Shop> = input.required<Shop>();

  dialogClosed: OutputEmitterRef<void> = output();
  shopUpdated: OutputEmitterRef<Shop> = output();

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
  ){}

  updateShop() {
    const sub = this.klcsAdminApi.updateShop(this.shop()).subscribe({
      next: val => this.shopUpdated.emit(val),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
