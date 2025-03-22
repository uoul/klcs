import { Component, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { Printer } from '../../domain/Printer';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap } from 'rxjs';
import { FormsModule } from '@angular/forms';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-create-printer-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './create-printer-dialog.component.html',
  styleUrl: './create-printer-dialog.component.css'
})
export class CreatePrinterDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<boolean> = output();
  printerCreated: OutputEmitterRef<Printer> = output();

  _printer: WritableSignal<Printer> = signal(new Printer())

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
    private notify: NotificationService,
  ){}

  createPrinter() {
    if(this._printer().Name.length > 0) {
      const sub = this.route.paramMap.pipe(
        mergeMap(params => this.shopAdminApi.createPrinterForShop(params.get("shopId") ?? "", this._printer()))
      ).subscribe({
        next: p => this.printerCreated.emit(p),
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
        complete: () => sub.unsubscribe(),
      })
    }
    this.init();
  }

  init(){
    this._printer.set(new Printer())
  }
}
