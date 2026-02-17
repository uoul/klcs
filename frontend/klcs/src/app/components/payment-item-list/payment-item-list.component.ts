import { Component, computed, input, InputSignal, model, ModelSignal, output, Signal, signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'klcs-payment-item-list',
  imports: [
    CommonModule,
  ],
  templateUrl: './payment-item-list.component.html',
  styleUrl: './payment-item-list.component.css'
})
export class PaymentItemListComponent {

  public articles: ModelSignal<Article[]> = model.required<Article[]>()
  public selected: WritableSignal<Article[]> = signal([])
  public sum: Signal<number> = computed(() => {
    let sum = 0;
    this.selected().forEach(a => {
      sum += a.Price
    })
    return sum
  })
  public closeDialog = output<void>()

  selectArticle(a: Article): void {
    if(!this.isSelected(a)){
      this.selected.set([a, ...this.selected()])
    } else {
      this.selected.set([...this.selected().filter(article => article !== a)])
    }
  }

  cashed(): void {
    this.articles.set(
      [...this.articles().filter(a => !!!this.selected().find(s => s === a))]
    )
    this.selected.set([])
    if(this.articles().length <= 0) {
      this.closeDialog.emit()
    }
  }

  handleMarkCheckbox(e: Event): void {
    const target = e.target as HTMLInputElement
    this.selected.set([])
    if(target.checked){
      this.selected.set([...this.articles()])
    }
  }

  isSelected(article: Article): boolean {
    return this.selected().find(a => article === a) ? true : false
  }
}
