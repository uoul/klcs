import { AfterViewInit, Component, input, InputSignal, model, OnInit, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import {ZXingScannerModule} from "@zxing/ngx-scanner";

@Component({
  selector: 'klcs-read-qr-dialog',
  imports: [
    ZXingScannerModule,
  ],
  templateUrl: './read-qr-dialog.component.html',
  styleUrl: './read-qr-dialog.component.css'
})
export class ReadQrDialogComponent implements AfterViewInit {
  dialogId: InputSignal<string> = input.required<string>();
  dialogClosed: OutputEmitterRef<void> = output();

  data = model<string>("");
  _dialog: HTMLDialogElement | null = null

  ngAfterViewInit(): void {
    this._dialog = document.getElementById(this.dialogId()) as HTMLDialogElement
  }

  onScanSuccess(data: string) {
    this.data.set(data)
    this._dialog?.close()
  }

  onScanError(error: Error){
    console.error(error)
    this._dialog?.close()
  }
}
