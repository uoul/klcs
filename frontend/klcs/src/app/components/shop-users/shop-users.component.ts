import { Component, computed, effect, input, OnInit, signal, untracked, WritableSignal } from '@angular/core';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ShopUser } from '../../domain/ShopUser';
import { ActivatedRoute } from '@angular/router';
import { mergeMap, take } from 'rxjs';
import { Role } from '../../domain/Role';
import { User } from '../../domain/User';
import { UserIdentity } from '../../domain/UserIdentity';
import { AuthService } from '../../services/auth/auth.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { CommonModule } from '@angular/common';
import { NotificationService } from '../../services/notification/notification.service';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';

interface _InternalUser {
  Id: string,
  Name: string,
  Username: string,
  roles: Map<Role, boolean>,
}

@Component({
  selector: 'klcs-shop-users',
  imports: [
    CommonModule,
    TranslatePipe,
  ],
  templateUrl: './shop-users.component.html',
  styleUrl: './shop-users.component.css'
})
export class ShopUsersComponent {

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private authService: AuthService,
    protected sellerApi: SellerApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){
    effect(() => {
      const shopId = this.sellerApi.getShopDetails().Id
      untracked(() => {
        this.refreshUsers(shopId);
        this.refreshRoles();
        this.currentUser.set(this.authService.getIdentity());
      })
    })
  }
  
  currentUser: WritableSignal<UserIdentity> = signal(new UserIdentity())
  klcsUsers: WritableSignal<ShopUser[]> = signal([])
  roles: WritableSignal<Role[]> = signal([])

  users = computed(() => {
    const res: _InternalUser[] = []
    for(let u of this.klcsUsers()){
      const userRoles: Map<Role, boolean> = new Map()
      for(let role of this.roles()) {
        userRoles.set(role, !!u.ShopRoles.find(r => r.Id === role.Id))
      }
      res.push({
        Id: u.Id,
        Name: u.Name,
        Username: u.Username,
        roles: userRoles
      })
    }
    return res
  })


  klcsConfig = KlcsConfig

  refreshRoles() {
    const sub = this.shopAdminApi.getRoles().subscribe({
      next: r => this.roles.set(r),
      error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
      complete: () => sub.unsubscribe(),
    })
  }

  refreshUsers(shopId: string) {
    const sub = this.shopAdminApi.getUsersForShop(shopId).subscribe({
      next: u => this.klcsUsers.set(u),
      error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
      complete: () => sub.unsubscribe(),
    })
  }

  userHasRole(user: ShopUser, role: Role): boolean {
    return user.ShopRoles.find(r => r.Id == role.Id) ? true : false
  }

  compareUser(a: _InternalUser, b: _InternalUser): number {
    if(a.Username == b.Username)
      return 0
    return a.Username < b.Username ? 1 : -1;
  }

  setUserRole(shopId: string, userId: string, role: Role, event: any) {
    if (event.target.checked) {
      const sub = this.shopAdminApi.addUserRoleForShop(shopId, userId, role).subscribe({
        next: _ => this.refreshUsers(shopId),
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => sub.unsubscribe(),
      })
    } else {
      const sub = this.shopAdminApi.deleteUserRoleForShop(shopId, userId, role.Id).subscribe({
        next: _ => this.refreshUsers(shopId),
        error: (err: HttpErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: this.translate.instant(`errors.${err.error?.Code}`)}),
        complete: () => sub.unsubscribe(),
      })
    }
  }
}
