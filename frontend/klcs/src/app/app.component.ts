import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AppBarComponent } from "./components/app-bar/app-bar.component";
import { SideNavComponent } from "./components/side-nav/side-nav.component";
import { AuthService } from './services/auth/auth.service';
import { CommonModule } from '@angular/common';
import { NotificationContainerComponent } from "./components/notification-container/notification-container.component";
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { LoadingService } from './services/loading/loading.service';

@Component({
  selector: 'klcs-root',
  imports: [
    RouterOutlet,
    AppBarComponent,
    SideNavComponent,
    CommonModule,
    NotificationContainerComponent,
    TranslatePipe,
],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  constructor(
    protected authService: AuthService,
    protected translate: TranslateService,
    protected loadingService: LoadingService,
  ){
    translate.addLangs(["de", "en"])
    translate.setFallbackLang("en")
    this.translate.use(this.translate.getBrowserLang() ?? 'en');
  }
}
