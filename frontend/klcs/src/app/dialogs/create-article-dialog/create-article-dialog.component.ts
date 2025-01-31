import { Component, computed, EventEmitter, Input, Output, Signal, signal, WritableSignal } from '@angular/core';
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

  @Input() dialogId!: string;
  @Input() printers: Signal<Printer[]> = signal([this.nonePrinter]);

  @Output() dialogClosed: EventEmitter<boolean> = new EventEmitter();
  @Output() articleCreated: EventEmitter<ArticleDetails> = new EventEmitter();

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
    const sub = this.route.paramMap.pipe(
      mergeMap(params => this.shopAdminApi.createArticle(params.get("shopId") ?? "", this.articleDetails))
    ).subscribe({
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
