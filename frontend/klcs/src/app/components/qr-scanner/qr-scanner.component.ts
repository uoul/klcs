import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  output,
  OutputEmitterRef,
  ViewChild,
} from '@angular/core';
import QrScanner from 'qr-scanner';

@Component({
  selector: 'klcs-qr-scanner',
  standalone: true,
  template: `<video #videoEl class="w-full rounded"></video>`,
  styles: [`
    :host { display: block; }
    video { width: 100%; border-radius: 0.25rem; }
  `],
})
export class QrScannerComponent implements AfterViewInit, OnDestroy {
  scanSuccess: OutputEmitterRef<string> = output<string>();
  scanError: OutputEmitterRef<Error> = output<Error>();

  @ViewChild('videoEl') private videoEl!: ElementRef<HTMLVideoElement>;
  private _scanner: QrScanner | null = null;

  ngAfterViewInit(): void {
    this._scanner = new QrScanner(
      this.videoEl.nativeElement,
      (result) => this.scanSuccess.emit(result.data),
      {
        preferredCamera: 'environment',
        highlightScanRegion: true,
        highlightCodeOutline: true,
        onDecodeError: (err) => {
          // qr-scanner fires this continuously while no code is in frame — filter noise
          if (err !== QrScanner.NO_QR_CODE_FOUND) {
            this.scanError.emit(err instanceof Error ? err : new Error(String(err)));
          }
        },
      }
    );
    this._scanner.start().catch((err) => this.scanError.emit(err));
  }

  ngOnDestroy(): void {
    this._scanner?.destroy();
  }
}
