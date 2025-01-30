import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppBarComponent } from "./components/app-bar/app-bar.component";
import { SideNavComponent } from "./components/side-nav/side-nav.component";
import { OAuthService } from 'angular-oauth2-oidc';
import { AuthService } from './services/auth/auth.service';
import { CommonModule } from '@angular/common';
import { KlcsConfig } from './config/KlcsConfig';

@Component({
  selector: 'klcs-root',
  imports: [
    RouterOutlet, 
    AppBarComponent, 
    SideNavComponent,
    CommonModule,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'klcs';

  constructor(
    private oauthService: OAuthService,
    protected authService: AuthService,
  ){
    this.oauthService.configure(KlcsConfig.OAuth);
    this.oauthService.setupAutomaticSilentRefresh();
    this.oauthService.loadDiscoveryDocumentAndLogin();
  }
}
