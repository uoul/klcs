import { Component, output } from '@angular/core';

@Component({
  selector: 'klcs-nav-item',
  imports: [],
  templateUrl: './nav-item.component.html',
  styleUrl: './nav-item.component.css'
})
export class NavItemComponent {
  clicked = output<void>();
}
