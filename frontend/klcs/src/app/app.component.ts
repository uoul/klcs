import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppBarComponent } from "./components/app-bar/app-bar.component";
import { SideNavComponent } from "./components/side-nav/side-nav.component";
import { AuthService } from './services/auth/auth.service';
import { CommonModule } from '@angular/common';
import { NotificationContainerComponent } from "./components/notification-container/notification-container.component";
import { TranslateService } from '@ngx-translate/core';

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
    protected translate: TranslateService,
  ){
    translate.addLangs(["de", "en", "it"])
    translate.setFallbackLang("en")
    this.translate.use(this.translate.getBrowserLang() ?? 'en');
  }
}
