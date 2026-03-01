import { Component, input, InputSignal, model, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { FormsModule } from '@angular/forms';
import { KlcsAdminApiService } from '../../services/klcs-admin-api/klcs-admin-api.service';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

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

  isActive: WritableSignal<boolean> = signal(false)

  constructor(
    private klcsAdminApi: KlcsAdminApiService,
    protected translate: TranslateService,
  ){}

  async updateShop() {
    if(!this.isActive()){
      this.isActive.set(true)
      try {
        const shop = await firstValueFrom(this.klcsAdminApi.updateShop(this.shop()))
        this.shop.set(shop)
      } finally { this.isActive.set(false) }
    }
  }
}
