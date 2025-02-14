import { Component, input } from '@angular/core';
import { Notification, NotificationService } from '../../services/notification/notification.service'


@Component({
  selector: 'klcs-notification-item',
  imports: [],
  templateUrl: './notification-item.component.html',
  styleUrl: './notification-item.component.css'
})
export class NotificationItemComponent {
  notification = input.required<Notification>()

  constructor(
    protected notificationService: NotificationService,
  ){}
}
