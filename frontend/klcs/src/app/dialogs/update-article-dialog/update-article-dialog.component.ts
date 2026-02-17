import { Component, computed, EventEmitter, input, Input, InputSignal, output, Output, OutputEmitterRef, signal, Signal, WritableSignal } from '@angular/core';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { Printer } from '../../domain/Printer';
import { FormsModule } from '@angular/forms';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

@Component({
  selector: 'klcs-update-article-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './update-article-dialog.component.html',
  styleUrl: './update-article-dialog.component.css'
})
export class UpdateArticleDialogComponent {
  nonePrinter: Printer = { Id: "", Name: "-", Connected: false };
  
  dialogId: InputSignal<string> = input.required<string>();
  printers: InputSignal<Printer[]> = input([this.nonePrinter]);
  article: InputSignal<ArticleDetails> = input(new ArticleDetails());

  dialogClosed: OutputEmitterRef<void> = output();
  articleUpdated: OutputEmitterRef<ArticleDetails> = output();

  _uiArticle = computed(() => {
    const details = this.article()
    details.Price = this.article().Price / 100
    details.Printer = this.article().Printer ?? this.nonePrinter
    return details
  })

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private notify: NotificationService,
  ){}

  updateArticle() {
    this.article().Price = Math.floor(this._uiArticle().Price * 100);
    this.article().Printer = (!this._uiArticle().Printer || this._uiArticle().Printer!.Id.length <= 0) ? null : this._uiArticle().Printer;
    const sub = this.shopAdminApi.updateArticle(this.article()).subscribe({
      next: a => this.articleUpdated.emit(a),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
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
