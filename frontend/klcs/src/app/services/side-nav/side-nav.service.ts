import { Injectable, Signal, signal, WritableSignal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class SideNavService {

  private readonly MOBILE_RESOLUTION = 1280;

  private _isOpen: WritableSignal<boolean> = signal(true);
  private _isMobile: WritableSignal<boolean> = signal(window.innerWidth < this.MOBILE_RESOLUTION);

  constructor() {
    addEventListener("resize", (event) => {
      const screenWidth = window.innerWidth
      if(this._isMobile() && screenWidth >= this.MOBILE_RESOLUTION)
        this._isMobile.set(false)
      else if(!this._isMobile() && screenWidth < this.MOBILE_RESOLUTION)
        this._isMobile.set(true)
    });
  }

  public get isOpen(): Signal<boolean> {
    return this._isOpen
  }

  public get isMobile(): Signal<boolean> {
    return this._isMobile
  }

  public toggle(): void {
    this._isOpen.set(!this.isOpen());
  }
}
