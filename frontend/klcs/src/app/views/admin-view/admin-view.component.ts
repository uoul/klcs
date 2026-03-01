import { Component, OnInit, signal, WritableSignal } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { CreateShopDialogComponent } from "../../dialogs/create-shop-dialog/create-shop-dialog.component";
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { UpdateShopDialogComponent } from "../../dialogs/update-shop-dialog/update-shop-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-admin-view',
  imports: [
    CreateShopDialogComponent,
    UpdateShopDialogComponent,
    TranslatePipe,
],
  templateUrl: './admin-view.component.html',
  styleUrl: './admin-view.component.css'
})
export class AdminViewComponent implements OnInit {
  
  readonly CREATE_DIALOG_ID = "create-dialog";
  readonly EDIT_DIALOG_ID = "create-edit";

  shops = signal<Shop[]>([]); 
  _currentSelectedShop = signal<Shop>(new Shop());
  
  constructor(
    private klcsAdminApi: KlcsAdminApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
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

  async deleteShop(shop: Shop) {
    if(confirm(this.translate.instant("views.admin.DeletePrompt", { name: shop.Name }))){
      try {
        await firstValueFrom(this.klcsAdminApi.deleteShop(shop.Id))
        this.notify.show({type: "success", duration: KlcsConfig.durationError, message: this.translate.instant("success.ShopDeleted")})
      } catch { }
    }
  }

  async refresh(){
    try {
      this.shops.set(
        await firstValueFrom(this.klcsAdminApi.getShops())
      )
    } catch {}
  }
}
