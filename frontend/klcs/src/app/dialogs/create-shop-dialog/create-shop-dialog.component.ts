import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Shop } from '../../domain/Shop';

@Component({
  selector: 'klcs-create-shop-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './create-shop-dialog.component.html',
  styleUrl: './create-shop-dialog.component.css'
})
export class CreateShopDialogComponent {
  @Input() dialogId!: string;
  @Output() dialogClosed: EventEmitter<boolean> = new EventEmitter();
  @Output() shopCreated: EventEmitter<Shop> = new EventEmitter();

  shopName: string = ""

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
  ){}

  createShop() {
    if(this.shopName.length > 0) {
      const shop: Shop = new Shop(undefined, this.shopName);
      const sub = this.klcsAdminApi.createShop(shop).subscribe({
        next: val => this.shopCreated.emit(val),
        error: err => console.error(err),
        complete: () => sub.unsubscribe()
      });
    }
    this.init();
  }

  init(){
    this.shopName = "";
  }

}
