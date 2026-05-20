import {
  AfterViewInit,
  ChangeDetectionStrategy,
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
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `<video #videoEl class="w-full rounded" playsinline></video>`,
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
  private _started = false;

  ngAfterViewInit(): void {
    this._scanner = new QrScanner(
      this.videoEl.nativeElement,
      (result) => this.scanSuccess.emit(result.data),
      {
        preferredCamera: 'environment',
        highlightScanRegion: true,
        highlightCodeOutline: true,
        onDecodeError: (err) => {
          if (err !== QrScanner.NO_QR_CODE_FOUND) {
            this.scanError.emit(err instanceof Error ? err : new Error(String(err)));
          }
        },
      }
    );

    this._scanner
      .start()
      .then(() => (this._started = true))
      .catch((err) => {
        // Do NOT emit scanError here if the parent uses it to destroy the component
        console.error('QrScanner failed to start:', err);
        // Only emit if you're sure the parent won't unmount on error
        // this.scanError.emit(err instanceof Error ? err : new Error(String(err)));
      });
  }

  ngOnDestroy(): void {
    if (this._scanner) {
      if (this._started) {
        this._scanner.stop();
      }
      this._scanner.destroy();
      this._scanner = null;
    }
  }
}
