import { Component, EventEmitter, Input, Output, signal, WritableSignal } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap } from 'rxjs';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'klcs-create-printer-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './create-printer-dialog.component.html',
  styleUrl: './create-printer-dialog.component.css'
})
export class CreatePrinterDialogComponent {
  @Input() dialogId!: string;
  @Output() dialogClosed: EventEmitter<boolean> = new EventEmitter();
  @Output() printerCreated: EventEmitter<Printer> = new EventEmitter();

  _printer: WritableSignal<Printer> = signal(new Printer())

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
  ){}

  createPrinter() {
    if(this._printer().Name.length > 0) {
      const sub = this.route.paramMap.pipe(
        mergeMap(params => this.shopAdminApi.createPrinterForShop(params.get("shopId") ?? "", this._printer()))
      ).subscribe({
        next: p => this.printerCreated.emit(p),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      })
    }
    this.init();
  }

  init(){
    this._printer.set(new Printer())
  }
}
