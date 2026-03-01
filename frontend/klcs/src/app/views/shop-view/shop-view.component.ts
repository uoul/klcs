import { Component, computed, OnInit, signal, WritableSignal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { SellerApiService } from '../../services/seller-api/seller-api.service';

import { CashdeskComponent } from "../../components/cashdesk/cashdesk.component";
import { KlcsConfig } from '../../config/KlcsConfig';
import { ShopArticlesComponent } from "../../components/shop-articles/shop-articles.component";
import { ShopPrintersComponent } from "../../components/shop-printers/shop-printers.component";
import { ShopUsersComponent } from "../../components/shop-users/shop-users.component";
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { finalize } from 'rxjs';

@Component({
  selector: 'klcs-shop-view',
  imports: [
    CashdeskComponent,
    ShopArticlesComponent,
    ShopPrintersComponent,
    ShopUsersComponent,
    TranslatePipe,
],
  templateUrl: './shop-view.component.html',
  styleUrl: './shop-view.component.css'
})
export class ShopViewComponent implements OnInit {

  currentTab = signal<number>(1)
  isAdmin = computed(() => this.sellerApi.getShopDetails().UserRoles.find(r => r == KlcsConfig.ShopRoleAdmin) ? true : false);

  constructor(
    protected sellerApi: SellerApiService,
    private route: ActivatedRoute,
    protected translate: TranslateService,
  ){}

  ngOnInit(): void {
    const sub = this.route.paramMap.pipe(
      finalize(() => sub.unsubscribe())
    ).subscribe({
      next: params => this.sellerApi.updateShopId(params.get("shopId") ?? ""),
    })
  }
}
