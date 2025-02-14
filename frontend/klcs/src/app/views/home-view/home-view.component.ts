import { Component } from '@angular/core';
import { Notification, NotificationService } from '../../services/notification/notification.service';

@Component({
  selector: 'klcs-home-view',
  imports: [],
  templateUrl: './home-view.component.html',
  styleUrl: './home-view.component.css'
})
export class HomeViewComponent {

  constructor(
    private notificationService: NotificationService,
  ){}

  createNotification(){
    const n: Notification = {
      duration: 3000,
      message: `${Math.random()}`,
      type: "warning"
    }
    this.notificationService.show(n)
  }
}
