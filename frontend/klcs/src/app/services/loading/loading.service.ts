import { computed, Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class LoadingService {
  private count = signal(0);
  isLoading = computed(() => this.count() > 0);

  increment() { this.count.update(n => n + 1); }
  decrement() { this.count.update(n => Math.max(0, n - 1)); }
}
