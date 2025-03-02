import { Component, effect, input, signal, untracked, WritableSignal } from '@angular/core';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { CreatePrinterDialogComponent } from "../../dialogs/create-printer-dialog/create-printer-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { SellerApiService } from '../../services/seller-api/seller-api.service';

@Component({
  selector: 'klcs-shop-printers',
  imports: [CreatePrinterDialogComponent],
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
  ){
    effect(() => {
      const shopId = this.sellerApi.getShopDetails().Id
      untracked(() => this.refreshPrinters()) 
    })
  }

  refreshPrinters() {
    const sub = this.shopAdminApi.getPrinters(this.sellerApi.getShopDetails().Id).subscribe({
      next: p =>  this._printers.set(p),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }

  deletePrinter(printer: Printer) {
    if(confirm(`Do you realy want to delete ${printer.Name}?`)){
      const sub = this.shopAdminApi.deletePrinter(printer.Id).subscribe({
        next: _ => this.refreshPrinters(),
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
        complete: () => sub.unsubscribe(),
      })
    }
  }

  showCreateDialog(){
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }
}
