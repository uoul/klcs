import { Component, effect, input, signal, untracked, WritableSignal } from '@angular/core';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { CreatePrinterDialogComponent } from "../../dialogs/create-printer-dialog/create-printer-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-shop-printers',
  imports: [
    CreatePrinterDialogComponent,
    TranslatePipe,
  ],
  templateUrl: './shop-printers.component.html',
  styleUrl: './shop-printers.component.css'
})
export class ShopPrintersComponent {
  protected readonly CREATE_DIALOG_ID = "create-printer-dialog"

  _printers: WritableSignal<Printer[]> = signal([])

  constructor(
    private shopAdminApi: ShopAdminApiService,
    protected sellerApi: SellerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){
    effect(() => {
      const shopId = this.sellerApi.getShopDetails().Id
      untracked(() => this.refreshPrinters()) 
    })
  }

  async refreshPrinters() {
    try {
      const printers = await firstValueFrom(this.shopAdminApi.getPrinters(this.sellerApi.getShopDetails().Id))
      this._printers.set(printers)
    } catch {}
  }

  async deletePrinter(printer: Printer) {
    if(confirm(this.translate.instant("components.shop-printers.DeletePrompt", { name: printer.Name }))){
      try {
        await firstValueFrom(this.shopAdminApi.deletePrinter(printer.Id))
        await this.refreshPrinters()
      } catch {}
    }
  }

  showCreateDialog(){
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }
}
