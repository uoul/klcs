import { Component, EventEmitter, Input, OnInit, Output, signal, WritableSignal } from '@angular/core';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap, take } from 'rxjs';
import { CreatePrinterDialogComponent } from "../../dialogs/create-printer-dialog/create-printer-dialog.component";

@Component({
  selector: 'klcs-shop-printers',
  imports: [CreatePrinterDialogComponent],
  templateUrl: './shop-printers.component.html',
  styleUrl: './shop-printers.component.css'
})
export class ShopPrintersComponent implements OnInit {
  protected readonly CREATE_DIALOG_ID = "create-printer-dialog"

  _printers: WritableSignal<Printer[]> = signal([])

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
  ){}

  ngOnInit(): void {
    this.refreshPrinters()
  }

  refreshPrinters() {
    const sub = this.route.paramMap.pipe(
      take(1),
      mergeMap(params => this.shopAdminApi.getPrinters(params.get("shopId") ?? "")),
    ).subscribe({
      next: p =>  this._printers.set(p),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  deletePrinter(printer: Printer) {
    if(confirm(`Do you realy want to delete ${printer.Name}?`)){
      const sub = this.shopAdminApi.deletePrinter(printer.Id).subscribe({
        next: _ => this.refreshPrinters(),
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
