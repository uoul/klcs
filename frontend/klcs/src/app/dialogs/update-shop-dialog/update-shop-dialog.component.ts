import { Component, EventEmitter, Input, OnInit, Output, Signal } from '@angular/core';
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
  @Input() dialogId!: string;
  @Input() shop!: Signal<Shop>;
  @Output() dialogClosed: EventEmitter<void> = new EventEmitter();
  @Output() shopUpdated: EventEmitter<Shop> = new EventEmitter();

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
  ){}

  updateShop() {
    const sub = this.klcsAdminApi.updateShop(this.shop()).subscribe({
      next: val => this.shopUpdated.emit(),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
