import { Component, OnInit, signal, WritableSignal } from '@angular/core';
import { SellerApiService } from '../../services/seller-api/seller-api.service';

import { HistoryItem } from '../../domain/HistoryItem';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { FormsModule } from '@angular/forms';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { finalize } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'klcs-history-view',
  imports: [FormsModule, TranslatePipe],
  templateUrl: './history-view.component.html',
  styleUrl: './history-view.component.css',
})
export class HistoryViewComponent implements OnInit {
  historyLength: WritableSignal<number> = signal(10);
  history: WritableSignal<HistoryItem[]> = signal([]);

  reprintRequestRunning: WritableSignal<boolean> = signal(false);

  constructor(
    private sellerApi: SellerApiService,
    private notify: NotificationService,
    private translate: TranslateService,
  ) {}

  ngOnInit(): void {
    this.refresh();
  }

  refresh(): void {
    const sub = this.sellerApi
      .getHistory(this.historyLength())
      .pipe(finalize(() => sub.unsubscribe()))
      .subscribe({
        next: (history) => this.history.set(history),
        error: (err) =>
          this.notify.show({
            type: 'error',
            duration: KlcsConfig.durationError,
            message: err,
          }),
      });
  }

  checkAnyNotPrinted(entry: HistoryItem): boolean {
    return entry.Articles.find((a) => !!!a.PrinterAck) ? true : false;
  }

  sendPrintJob(transactionId: string) {
    if (!this.reprintRequestRunning()) {
      this.reprintRequestRunning.set(true);
      const sub = this.sellerApi
        .reprintOrder(transactionId)
        .pipe(finalize(() => {this.reprintRequestRunning.set(false); sub.unsubscribe()}))
        .subscribe({
          next: (_) => setTimeout(() => this.refresh(), 1000),
          error: (err: HttpErrorResponse) =>
            this.notify.show({
              type: 'error',
              duration: KlcsConfig.durationError,
              message: this.translate.instant(`errors.${err.error?.Code}`),
            }),
        });
    }
  }
}
