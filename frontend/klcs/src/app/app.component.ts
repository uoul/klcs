import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppBarComponent } from "./components/app-bar/app-bar.component";
import { SideNavComponent } from "./components/side-nav/side-nav.component";
import { OAuthService } from 'angular-oauth2-oidc';
import { authCodeFlowConfig } from './config/authConfig';

@Component({
  selector: 'klcs-root',
  imports: [
    RouterOutlet, 
    AppBarComponent, 
    SideNavComponent,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'klcs';

  constructor(
    private oauthService: OAuthService,
  ){
    this.oauthService.configure(authCodeFlowConfig);
    this.oauthService.setupAutomaticSilentRefresh();
    this.oauthService.loadDiscoveryDocumentAndLogin();
  }
}
