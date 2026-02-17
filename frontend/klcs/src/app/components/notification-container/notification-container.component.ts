import { Component } from '@angular/core';
import {  NotificationService } from '../../services/notification/notification.service';
import { NotificationItemComponent } from '../notification-item/notification-item.component';
import { trigger, transition, style, animate } from '@angular/animations';

@Component({
  selector: 'klcs-notification-container',
  imports: [
    NotificationItemComponent,
  ],
  templateUrl: './notification-container.component.html',
  styleUrl: './notification-container.component.css',
  animations: [
    trigger('toastAnimation', [
      // Fade-in und Slide-up Animation
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(10px)' }),
        animate('300ms ease-out', style({ opacity: 1, transform: 'translateY(0)' }))
      ]),
    ])
  ]
})
export class NotificationContainerComponent {

  constructor(
    protected notificationService: NotificationService,
  ){}

}
