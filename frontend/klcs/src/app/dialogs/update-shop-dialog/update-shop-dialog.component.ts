import { Component, EventEmitter, input, Input, InputSignal, model, OnInit, output, Output, OutputEmitterRef, Signal } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { FormsModule } from '@angular/forms';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'klcs-update-shop-dialog',
  imports: [
    FormsModule,
    TranslatePipe,
  ],
  templateUrl: './update-shop-dialog.component.html',
  styleUrl: './update-shop-dialog.component.css'
})
export class UpdateShopDialogComponent {
  dialogId: InputSignal<string> = input.required<string>();
  shop = model.required<Shop>();

  dialogClosed: OutputEmitterRef<void> = output();

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  updateShop() {
    const sub = this.klcsAdminApi.updateShop(this.shop()).subscribe({
      next: val => this.shop.set(val),
      error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
      complete: () => sub.unsubscribe(),
    })
  }
}
