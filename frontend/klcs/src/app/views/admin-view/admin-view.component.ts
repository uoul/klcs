import { Component, OnInit, signal, WritableSignal } from '@angular/core';
import { Observable, subscribeOn } from 'rxjs';
import { Shop } from '../../domain/Shop';
import { CommonModule } from '@angular/common';
import { CreateShopDialogComponent } from "../../dialogs/create-shop-dialog/create-shop-dialog.component";
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { UpdateShopDialogComponent } from "../../dialogs/update-shop-dialog/update-shop-dialog.component";

@Component({
  selector: 'klcs-admin-view',
  imports: [
    CommonModule,
    CreateShopDialogComponent,
    UpdateShopDialogComponent
],
  templateUrl: './admin-view.component.html',
  styleUrl: './admin-view.component.css'
})
export class AdminViewComponent implements OnInit {
  
  readonly CREATE_DIALOG_ID = "create-dialog";
  readonly EDIT_DIALOG_ID = "create-edit";

  shops: Observable<Shop[]> = new Observable(); 
  _currentSelectedShop: WritableSignal<Shop> = signal(new Shop());

  J = JSON;
  
  constructor(
    private klcsAdminApi: KlcsAdminApiService,
  ){}

  ngOnInit(): void {
    this.refresh();
  }

  showCreateDialog() {
    const dialog = (document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement)
    dialog.showModal();
  }

  showEditDialog(shop: Shop) {
    const dialog = (document.getElementById(this.EDIT_DIALOG_ID) as HTMLDialogElement)
    this._currentSelectedShop.set(JSON.parse(JSON.stringify(shop)))
    dialog.showModal();
  }

  deleteShop(shop: Shop) {
    if(confirm(`Do you realy want to delete shop ${shop.Name}?`)){
      const sub = this.klcsAdminApi.deleteShop(shop.Id).subscribe({
        next: _ => console.log("deleted successfully"),
        error: err => console.error(err),
        complete: () => {
          this.refresh();
          sub.unsubscribe();
        },
      });
    }
  }

  refresh(){
    this.shops = this.klcsAdminApi.getShops();
  }
}
