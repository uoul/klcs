import { Component, computed, EventEmitter, input, Input, InputSignal, output, Output, OutputEmitterRef, Signal, signal, WritableSignal } from '@angular/core';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { mergeMap } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Printer } from '../../domain/Printer';

@Component({
  selector: 'klcs-create-article-dialog',
  imports: [
    FormsModule,
  ],
  templateUrl: './create-article-dialog.component.html',
  styleUrl: './create-article-dialog.component.css'
})
export class CreateArticleDialogComponent {
  nonePrinter: Printer = { Id: "", Name: "-" };

  dialogId: InputSignal<string> = input.required<string>();
  shopId: InputSignal<string> = input.required<string>();
  printers: InputSignal<Printer[]> = input([this.nonePrinter]);

  dialogClosed: OutputEmitterRef<void> = output();
  articleCreated: OutputEmitterRef<ArticleDetails> = output();

  articleDetails: ArticleDetails = new ArticleDetails()
  priceUi: number = 0.0;
  printer: Printer = this.nonePrinter;

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
  ){}

  createArticle() {
    this.articleDetails.Price = Math.floor(this.priceUi * 100);
    this.articleDetails.Printer = this.printer.Id === "" ? null : this.printer;
    const sub = this.shopAdminApi.createArticle(this.shopId(), this.articleDetails).subscribe({
      next: a => this.articleCreated.emit(a),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
    this.init()
  }

  init(){
    this.articleDetails = new ArticleDetails()
    this.priceUi = 0.0
    this.printer = this.nonePrinter;
  }
}
