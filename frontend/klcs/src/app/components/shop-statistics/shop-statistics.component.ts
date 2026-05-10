import { Component, computed, effect, ElementRef, OnDestroy, Signal, signal, ViewChild, WritableSignal } from '@angular/core';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { CommonModule } from '@angular/common';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { firstValueFrom } from 'rxjs';
import { RevenueItem } from '../../domain/RevenueItem';
import { Chart, registerables } from 'chart.js';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';

Chart.register(...registerables);

@Component({
  selector: 'klcs-shop-statistics',
  imports: [CommonModule, TranslatePipe],
  templateUrl: './shop-statistics.component.html',
  styleUrl: './shop-statistics.component.css',
})
export class ShopStatisticsComponent implements OnDestroy {

  revenue: WritableSignal<RevenueItem[] | null> = signal(null);
  stock: Signal<{ Article: string; Amount: number }[] | null> = computed(() => {
    const articles = Object.values(this.sellerApi.getShopDetails().Categories).flat();
    return articles
      .filter(a => a.StockAmount)
      .map<{ Article: string; Amount: number }>(a => ({ Article: a.Name, Amount: a.StockAmount! }));
  });


  @ViewChild('sellsChart', { static: true }) private sellsChartRef!: ElementRef<HTMLCanvasElement>;
  @ViewChild('stockChart', { static: true }) private stockChartRef!: ElementRef<HTMLCanvasElement>;

  private sellsChart!: Chart;
  private stockChart!: Chart;

  totalRevenue: Signal<number> = computed(() =>
    (this.revenue() ?? []).reduce((sum, item) => sum + item.Sum, 0)
  );

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private sellerApi: SellerApiService,
    protected translate: TranslateService,
  ) {
    effect(() => {
      const shopId = this.sellerApi.getShopDetails().Id;
      this.refresh(shopId);
    });

    effect(() => {
      const data = this.revenue();
      if (!data) return;
      this.updateSellsChart(data);
    });

    effect(() => {
      const data = this.stock();
      if (!data) return;
      this.updateStockChart(data);
    });
  }

  ngOnDestroy() {
    this.sellsChart?.destroy();
    this.stockChart?.destroy();
  }

  private updateSellsChart(data: RevenueItem[]): void {
    const labels = data.map(item => item.Article);
    const amounts = data.map(item => item.Amount);
    const sums = data.map(item => item.Sum);

    if (this.sellsChart) {
      this.sellsChart.data.labels = labels;
      this.sellsChart.data.datasets[0].data = amounts;
      this.sellsChart.data.datasets[1].data = sums;
      this.sellsChart.update();
    } else {
      this.sellsChart = new Chart(this.sellsChartRef.nativeElement, {
        type: 'bar',
        data: {
          labels,
          datasets: [
            {
              label: this.translate.instant('components.shop-statistics.ChartSellsLblAmount'),
              data: amounts,
              backgroundColor: 'rgba(230, 120, 23, 0.7)',
              borderColor: 'rgba(186, 96, 18, 1)',
              borderWidth: 1.5,
            },
            {
              label: this.translate.instant('components.shop-statistics.ChartSellsLblRevenue'),
              data: sums,
              backgroundColor: 'rgba(48, 48, 48, 0.7)',
              borderColor: 'rgba(32, 32, 32, 1)',
              borderWidth: 1.5,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { position: 'top' } },
          scales: { y: { beginAtZero: true } },
          indexAxis: 'y',
        },
      });
    }
  }

  private updateStockChart(data: { Article: string; Amount: number }[]): void {
    const labels = data.map(item => item.Article);
    const amounts = data.map(item => item.Amount);

    if (this.stockChart) {
      this.stockChart.data.labels = labels;
      this.stockChart.data.datasets[0].data = amounts;
      this.stockChart.update();
    } else {
      this.stockChart = new Chart(this.stockChartRef.nativeElement, {
        type: 'bar',
        data: {
          labels,
          datasets: [
            {
              label: this.translate.instant('components.shop-statistics.ChartStockLblAmount'),
              data: amounts,
              backgroundColor: 'rgba(230, 120, 23, 0.7)',
              borderColor: 'rgba(186, 96, 18, 1)',
              borderWidth: 1.5,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { position: 'top' } },
          scales: { y: { beginAtZero: true } },
          indexAxis: 'y',
        },
      });
    }
  }

  public async refresh(shopId?: string): Promise<void> {
    this.revenue.set(null);
    try {
      const details = this.sellerApi.getShopDetails();
      const id = shopId ?? details.Id;
      this.revenue.set(await firstValueFrom(this.shopAdminApi.getRevenue(id)));
    } catch (_) {}
  }
}
