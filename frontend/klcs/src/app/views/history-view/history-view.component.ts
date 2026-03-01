import { Component, OnInit, signal, WritableSignal } from '@angular/core';
import { SellerApiService } from '../../services/seller-api/seller-api.service';

import { HistoryItem } from '../../domain/HistoryItem';
import { FormsModule } from '@angular/forms';
import { TranslatePipe } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

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
  refreshActive: WritableSignal<boolean> = signal(false)

  constructor(
    private sellerApi: SellerApiService,
  ) {}

  ngOnInit(): void {
    this.refresh();
  }

  async refresh(): Promise<void> {
    if(!this.refreshActive()){
      this.refreshActive.set(true)
      try {
        const history = await firstValueFrom(this.sellerApi.getHistory(this.historyLength()))
        this.history.set(history)
      } finally { this.refreshActive.set(false) }
    }
   
  }

  checkAnyNotPrinted(entry: HistoryItem): boolean {
    return entry.Articles.find((a) => !!!a.PrinterAck) ? true : false;
  }

  async sendPrintJob(transactionId: string) {
    if (!this.reprintRequestRunning()) {
      this.reprintRequestRunning.set(true)
      try {
        await firstValueFrom(this.sellerApi.reprintOrder(transactionId))
        const t = setTimeout(() => { clearTimeout(t); this.refresh() }, 1000)
      } finally {
        this.reprintRequestRunning.set(false)
      }
    }
  }
}
