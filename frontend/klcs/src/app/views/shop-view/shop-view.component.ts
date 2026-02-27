import { Component, computed, OnInit, signal, WritableSignal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { ShopDetails } from '../../domain/ShopDetails';

import { CashdeskComponent } from "../../components/cashdesk/cashdesk.component";
import { KlcsConfig } from '../../config/KlcsConfig';
import { ShopArticlesComponent } from "../../components/shop-articles/shop-articles.component";
import { ShopPrintersComponent } from "../../components/shop-printers/shop-printers.component";
import { ShopUsersComponent } from "../../components/shop-users/shop-users.component";
import { NotificationService } from '../../services/notification/notification.service';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'klcs-shop-view',
  imports: [
    CashdeskComponent,
    ShopArticlesComponent,
    ShopPrintersComponent,
    ShopUsersComponent
],
  templateUrl: './shop-view.component.html',
  styleUrl: './shop-view.component.css'
})
export class ShopViewComponent implements OnInit {

  currentTab = signal<number>(1)
  isAdmin = computed(() => this.sellerApi.getShopDetails().UserRoles.find(r => r == KlcsConfig.ShopRoleAdmin) ? true : false);

  constructor(
    protected sellerApi: SellerApiService,
    private notify: NotificationService,
    private route: ActivatedRoute,
    protected translate: TranslateService,
  ){}

  ngOnInit(): void {
    const sub = this.route.paramMap.subscribe({
      next: params => this.sellerApi.updateShopId(params.get("shopId") ?? ""),
      error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
      complete: () => sub.unsubscribe()
    })
  }
}
