import { Component, computed, EventEmitter, Input, Output, signal, Signal, WritableSignal } from '@angular/core';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { Printer } from '../../domain/Printer';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'klcs-update-article-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './update-article-dialog.component.html',
  styleUrl: './update-article-dialog.component.css'
})
export class UpdateArticleDialogComponent {
  nonePrinter: Printer = { Id: "", Name: "-" };
  
  @Input() dialogId!: string;
  @Input() printers: Signal<Printer[]> = signal([this.nonePrinter]);
  @Input() article: Signal<ArticleDetails> = signal(new ArticleDetails());

  @Output() dialogClosed: EventEmitter<boolean> = new EventEmitter();
  @Output() articleUpdated: EventEmitter<ArticleDetails> = new EventEmitter();

  _uiArticle = computed(() => {
    const details = this.article()
    details.Price = this.article().Price / 100
    details.Printer = this.article().Printer ?? this.nonePrinter
    return details
  })

  constructor(
    private shopAdminApi: ShopAdminApiService,
  ){}

  updateArticle() {
    this.article().Price = Math.floor(this._uiArticle().Price * 100);
    this.article().Printer = (!this._uiArticle().Printer || this._uiArticle().Printer!.Id.length <= 0) ? null : this._uiArticle().Printer;
    const sub = this.shopAdminApi.updateArticle(this.article()).subscribe({
      next: a => this.articleUpdated.emit(a),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  comparePriner(p1: Printer, p2: Printer) {
    if(!p1)
      return false;
    if(!p2)
      return false;
    return p1.Id == p2.Id;
  }
}
