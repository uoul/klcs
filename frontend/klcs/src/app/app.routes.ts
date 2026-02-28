import { Routes } from '@angular/router';
import { ShopViewComponent } from './views/shop-view/shop-view.component';
import { AccountViewComponent } from './views/account-view/account-view.component';
import { AdminViewComponent } from './views/admin-view/admin-view.component';
import { authGuard } from './guards/auth.guard';
import { HistoryViewComponent } from './views/history-view/history-view.component';

export const routes: Routes = [
  {
    path: "",
    pathMatch: "full",
    redirectTo: "home"
  },
  {
    path: "shops/:shopId",
    canActivate: [ authGuard ],
    component: ShopViewComponent,
  },
  {
    path: "history",
    canActivate: [ authGuard ],
    component: HistoryViewComponent,
  },
  {
    path: "accounts",
    canActivate: [ authGuard ],
    component: AccountViewComponent,
  },
  {
    path: "admin",
    canActivate: [ authGuard ],
    component: AdminViewComponent,
  },
];
