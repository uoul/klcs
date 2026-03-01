import { AfterViewInit, Component, OnInit, signal, WritableSignal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavItemComponent } from "../nav-item/nav-item.component";
import { SideNavService } from '../../services/side-nav/side-nav.service';
import { AuthService } from '../../services/auth/auth.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { SellerApiService } from '../../services/seller-api/seller-api.service';

import { Shop } from '../../domain/Shop';
import { CheckAccountDialogComponent } from "../../dialogs/check-account-dialog/check-account-dialog.component";
import { ChargeAccountDialogComponent } from "../../dialogs/charge-account-dialog/charge-account-dialog.component";
import { CloseAccountDialogComponent } from "../../dialogs/close-account-dialog/close-account-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { PublicApiService } from '../../services/public-api/public-api.service';
import { TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import {TranslatePipe } from "@ngx-translate/core";
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-side-nav',
  imports: [
    RouterModule,
    NavItemComponent,
    CheckAccountDialogComponent,
    ChargeAccountDialogComponent,
    CloseAccountDialogComponent,
    TranslatePipe,
],
  templateUrl: './side-nav.component.html',
  styleUrl: './side-nav.component.css'
})
export class SideNavComponent implements OnInit, AfterViewInit {
  constructor(
    protected sideNav: SideNavService,
    protected authService: AuthService,
    protected sellerApi: SellerApiService,
    private notify: NotificationService,
    protected publicApi: PublicApiService,
    protected translate: TranslateService,
  ){}

  protected readonly CHECK_CREDIT_DIALOG_ID = "check-credit-dialog"
  protected readonly CHARGE_CREDIT_DIALOG_ID = "charge-credit-dialog"
  protected readonly CLOSE_ACCOUNT_DIALOG_ID = "close-account-dialog"

  _checkAccountDialog: HTMLDialogElement|null = null
  _chargeAccountDialog: HTMLDialogElement|null = null
  _closeAccountDialog: HTMLDialogElement|null = null

  shops: WritableSignal<Shop[]> = signal([]);
  klcsConfig = KlcsConfig;

  async ngOnInit(): Promise<void> {
    try {
      const shops = await firstValueFrom(this.sellerApi.getShops())
      this.shops.set(shops)
    } catch {}
  }

  ngAfterViewInit(): void {
    this._checkAccountDialog = document.getElementById(this.CHECK_CREDIT_DIALOG_ID) as HTMLDialogElement
    this._chargeAccountDialog = document.getElementById(this.CHARGE_CREDIT_DIALOG_ID) as HTMLDialogElement
    this._closeAccountDialog = document.getElementById(this.CLOSE_ACCOUNT_DIALOG_ID) as HTMLDialogElement
  }

  checkUserRole(role: string): boolean {
    return this.authService.getIdentity().roles.find((r) => r == role) ? true : false;
  }

  updateMenuState() {
    if(this.sideNav.isMobile() && this.sideNav.isOpen())
      this.sideNav.toggle()
  }
}
