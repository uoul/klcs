import { Component, input, InputSignal, model, ModelSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Shop } from '../../domain/Shop';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

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

  isActive: WritableSignal<boolean> = signal(false)

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  async createShop() {
    if(!this.isActive()){
      this.isActive.set(true)
      try {
        const created = await firstValueFrom(this.klcsAdminApi.createShop(this._shop()))
        this.shopCreated.emit(created)
      } catch {}
      finally { this.init(); this.isActive.set(false) }
    }
  }

  init(){
    this._shop.set(new Shop(undefined, ""))
  }
}
