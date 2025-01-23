import { Injectable, Signal, signal, WritableSignal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class SideNavService {

  private _isOpen: WritableSignal<boolean> = signal(true);

  constructor() { }

  public get isOpen(): Signal<Boolean> {
    return this._isOpen
  }

  public toggle(): void {
    this._isOpen.set(!this.isOpen());
  }
}
