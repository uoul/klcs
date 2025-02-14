import { AfterViewInit, Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavItemComponent } from "../nav-item/nav-item.component";
import { SideNavService } from '../../services/side-nav/side-nav.service';
import { AuthService } from '../../services/auth/auth.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { CommonModule } from '@angular/common';
import { Shop } from '../../domain/Shop';
import { subscribeOn } from 'rxjs';
import { CheckAccountDialogComponent } from "../../dialogs/check-account-dialog/check-account-dialog.component";
import { ChargeAccountDialogComponent } from "../../dialogs/charge-account-dialog/charge-account-dialog.component";
import { CloseAccountDialogComponent } from "../../dialogs/close-account-dialog/close-account-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-side-nav',
  imports: [
    CommonModule,
    RouterModule,
    NavItemComponent,
    CheckAccountDialogComponent,
    ChargeAccountDialogComponent,
    CloseAccountDialogComponent
],
  templateUrl: './side-nav.component.html',
  styleUrl: './side-nav.component.css'
})
export class SideNavComponent implements OnInit, AfterViewInit {
  constructor(
    protected sideNav: SideNavService,
    protected authService: AuthService,
    protected sellerApi: SellerApiService,
    private notify: NotificationService
  ){}

  protected readonly CHECK_CREDIT_DIALOG_ID = "check-credit-dialog"
  protected readonly CHARGE_CREDIT_DIALOG_ID = "charge-credit-dialog"
  protected readonly CLOSE_ACCOUNT_DIALOG_ID = "close-account-dialog"

  _checkAccountDialog: HTMLDialogElement|null = null
  _chargeAccountDialog: HTMLDialogElement|null = null
  _closeAccountDialog: HTMLDialogElement|null = null

  shops: Shop[] = [];
  klcsConfig = KlcsConfig;

  ngOnInit(): void {
    const sub = this.sellerApi.getShops().subscribe({
      next: val => this.shops = val,
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationMedium, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
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
