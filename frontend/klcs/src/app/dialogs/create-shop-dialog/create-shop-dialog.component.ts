import { Component, input, InputSignal, model, ModelSignal, output, OutputEmitterRef } from '@angular/core';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Shop } from '../../domain/Shop';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'klcs-create-shop-dialog',
  imports: [
    FormsModule,
    TranslatePipe,
  ],
  templateUrl: './create-shop-dialog.component.html',
  styleUrl: './create-shop-dialog.component.css'
})
export class CreateShopDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output()
  shopCreated: OutputEmitterRef<Shop> = output()

  _shop: ModelSignal<Shop> = model(new Shop(undefined, ""))

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  createShop() {
    if(this._shop().Name.length > 0) {
      const sub = this.klcsAdminApi.createShop(this._shop()).subscribe({
        next: val => this.shopCreated.emit(val),
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => sub.unsubscribe()
      });
    }
    this.init();
  }

  init(){
    this._shop.set(new Shop(undefined, ""))
  }
}
