import { Component, computed, OnInit, signal, WritableSignal } from '@angular/core';
import { ActivatedRoute, Route } from '@angular/router';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { mergeMap, Observable } from 'rxjs';
import { ShopDetails } from '../../domain/ShopDetails';
import { CommonModule } from '@angular/common';
import { CashdeskComponent } from "../../components/cashdesk/cashdesk.component";
import { KlcsConfig } from '../../config/KlcsConfig';
import { ShopArticlesComponent } from "../../components/shop-articles/shop-articles.component";
import { ShopPrintersComponent } from "../../components/shop-printers/shop-printers.component";
import { ShopUsersComponent } from "../../components/shop-users/shop-users.component";
import { Article } from '../../domain/Article';

@Component({
  selector: 'klcs-shop-view',
  imports: [
    CommonModule,
    CashdeskComponent,
    ShopArticlesComponent,
    ShopPrintersComponent,
    ShopUsersComponent
],
  templateUrl: './shop-view.component.html',
  styleUrl: './shop-view.component.css'
})
export class ShopViewComponent implements OnInit {

  shop: WritableSignal<ShopDetails> = signal(new ShopDetails());

  articles = computed(() => this.shop().Categories);
  isAdmin = computed(() => this.shop().UserRoles.find(r => r == KlcsConfig.ShopRoleAdmin) ? true : false);

  constructor(
    private route: ActivatedRoute,
    private sellerApi: SellerApiService,
  ){}

  ngOnInit(): void {
    this.refresh();
  }

  refresh(): void {
    const sub = this.route.paramMap.pipe(
      mergeMap(params => this.sellerApi.getShopDetails(params.get("shopId") ?? ""))
    ).subscribe({
      next: resp => this.shop.set(resp),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
