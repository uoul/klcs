import { Injectable, Signal, signal, WritableSignal } from '@angular/core';

export type Notification = {
  type: 'info' | 'warning' | 'success' | 'error'
  duration: number
  message: string
}


@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  _notifications: WritableSignal<Notification[]> = signal([])

  constructor() { }

  show(msg: Notification) {
    this._notifications.update(current => [...current, msg])
    if(msg.duration > 0) {
      setTimeout(() => this.remove(msg), msg.duration)
    }
  }

  remove(msg: Notification) {
    this._notifications.update(current => current.filter(n => n != msg))
  }

  public get notifications(): Signal<Notification[]> {
    return this._notifications
  }
}
