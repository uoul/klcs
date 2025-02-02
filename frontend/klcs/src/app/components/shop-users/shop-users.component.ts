import { Component, computed, OnInit, signal, WritableSignal } from '@angular/core';
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
  ],
  templateUrl: './shop-users.component.html',
  styleUrl: './shop-users.component.css'
})
export class ShopUsersComponent implements OnInit {

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
    private authService: AuthService,
  ){}

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

  ngOnInit(): void {
    this.refreshUsers();
    this.refreshRoles();
    this.currentUser.set(this.authService.getIdentity());
  }

  refreshRoles() {
    const sub = this.shopAdminApi.getRoles().subscribe({
      next: r => this.roles.set(r),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  refreshUsers() {
    const sub = this.route.paramMap.pipe(
      mergeMap(params => this.shopAdminApi.getUsersForShop(params.get("shopId") ?? ""))
    ).subscribe({
      next: u => this.klcsUsers.set(u),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  userHasRole(user: ShopUser, role: Role): boolean {
    return user.ShopRoles.find(r => r.Id == role.Id) ? true : false
  }

  setUserRole(userId: string, role: Role, event: any) {
    if (event.target.checked) {
      const sub = this.route.paramMap.pipe(
        take(1),
        mergeMap(params => this.shopAdminApi.addUserRoleForShop(params.get("shopId") ?? "", userId, role))
      ).subscribe({
        next: _ => this.refreshUsers(),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      })
    } else {
      const sub = this.route.paramMap.pipe(
        take(1),
        mergeMap(params => this.shopAdminApi.deleteUserRoleForShop(params.get("shopId") ?? "", userId, role.Id))
      ).subscribe({
        next: _ => this.refreshUsers(),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      })
    }
  }
}
