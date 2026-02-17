import { Component, OnInit, Signal, signal } from '@angular/core';
import { SideNavService } from '../../services/side-nav/side-nav.service';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'klcs-app-bar',
  imports: [],
  templateUrl: './app-bar.component.html',
  styleUrl: './app-bar.component.css'
})
export class AppBarComponent {

  constructor(
    protected sideNav: SideNavService,
    protected authService: AuthService,
  ){}
  
  toggleMenu() {
    this.sideNav.toggle();
  }
}
