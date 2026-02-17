import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppBarComponent } from "./components/app-bar/app-bar.component";
import { SideNavComponent } from "./components/side-nav/side-nav.component";
import { OAuthService } from 'angular-oauth2-oidc';
import { AuthService } from './services/auth/auth.service';
import { CommonModule } from '@angular/common';
import { KlcsConfig } from './config/KlcsConfig';
import { NotificationContainerComponent } from "./components/notification-container/notification-container.component";
import { PublicApiService } from './services/public-api/public-api.service';

@Component({
  selector: 'klcs-root',
  imports: [
    RouterOutlet,
    AppBarComponent,
    SideNavComponent,
    CommonModule,
    NotificationContainerComponent
],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  constructor(
    protected authService: AuthService,
  ){}
}
