import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavItemComponent } from "../nav-item/nav-item.component";
import { SideNavService } from '../../services/side-nav/side-nav.service';
import { AuthService } from '../../services/auth/auth.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { CommonModule } from '@angular/common';
import { Shop } from '../../domain/Shop';
import { subscribeOn } from 'rxjs';

@Component({
  selector: 'klcs-side-nav',
  imports: [
    CommonModule,
    RouterModule,
    NavItemComponent
],
  templateUrl: './side-nav.component.html',
  styleUrl: './side-nav.component.css'
})
export class SideNavComponent implements OnInit {
  constructor(
    protected sideNav: SideNavService,
    protected authService: AuthService,
    protected sellerApi: SellerApiService,
  ){}

  shops: Shop[] = [];
  klcsConfig = KlcsConfig;

  ngOnInit(): void {
    const sub = this.sellerApi.getShops().subscribe({
      next: val => this.shops = val,
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  checkUserRole(role: string): boolean {
    return this.authService.getIdentity().roles.find((r) => r == role) ? true : false;
  }
}
