import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavItemComponent } from "../nav-item/nav-item.component";
import { SideNavService } from '../../services/side-nav/side-nav.service';

@Component({
  selector: 'klcs-side-nav',
  imports: [
    RouterModule,
    NavItemComponent
],
  templateUrl: './side-nav.component.html',
  styleUrl: './side-nav.component.css'
})
export class SideNavComponent {
  constructor(
    protected sideNav: SideNavService,
  ){}
}
