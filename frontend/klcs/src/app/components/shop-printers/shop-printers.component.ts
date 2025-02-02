import { Component, effect, input, signal, untracked, WritableSignal } from '@angular/core';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { CreatePrinterDialogComponent } from "../../dialogs/create-printer-dialog/create-printer-dialog.component";

@Component({
  selector: 'klcs-shop-printers',
  imports: [CreatePrinterDialogComponent],
  templateUrl: './shop-printers.component.html',
  styleUrl: './shop-printers.component.css'
})
export class ShopPrintersComponent {
  protected readonly CREATE_DIALOG_ID = "create-printer-dialog"

  shopId = input.required<string>()
  _printers: WritableSignal<Printer[]> = signal([])

  constructor(
    private shopAdminApi: ShopAdminApiService,
  ){
    effect(() => {
      const shopId = this.shopId()
      untracked(() => this.refreshPrinters(shopId)) 
    })
  }

  refreshPrinters(shopId: string) {
    const sub = this.shopAdminApi.getPrinters(shopId).subscribe({
      next: p =>  this._printers.set(p),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  deletePrinter(printer: Printer) {
    if(confirm(`Do you realy want to delete ${printer.Name}?`)){
      const sub = this.shopAdminApi.deletePrinter(printer.Id).subscribe({
        next: _ => this.refreshPrinters(this.shopId()),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      })
    }
  }

  showCreateDialog(){
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }
}
