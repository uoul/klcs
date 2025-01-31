import { Component, EventEmitter, Input, OnInit, Output, signal, Signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap, subscribeOn } from 'rxjs';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { CreateArticleDialogComponent } from "../../dialogs/create-article-dialog/create-article-dialog.component";
import { Printer } from '../../domain/Printer';

@Component({
  selector: 'klcs-shop-articles',
  imports: [
    CommonModule,
    CreateArticleDialogComponent
],
  templateUrl: './shop-articles.component.html',
  styleUrl: './shop-articles.component.css'
})
export class ShopArticlesComponent implements OnInit {
  @Input() categories: Signal<Map<string, Article[]>> = signal(new Map());
  @Output() articlesChanged: EventEmitter<void> = new EventEmitter();

  protected readonly CREATE_DIALOG_ID = "create-article-dialog"

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
  ){}

  _printers: WritableSignal<Printer[]> = signal([])

  ngOnInit(): void {
    this.refreshPrinters()
  }

  deleteArticle(articleId: string) {
    if(confirm(`Do you realy want to delete Article?`)){
      const sub = this.shopAdminApi.deleteArticle(articleId).subscribe({
        next: _ => this.articlesChanged.emit(),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      });
    }
  }

  updateArticle(article: ArticleDetails) {
    const sub = this.shopAdminApi.updateArticle(article).subscribe({
      next: _ => this.articlesChanged.emit(),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  refreshPrinters() {
    const sub = this.route.paramMap.pipe(
      mergeMap(params => this.shopAdminApi.getPrinters(params.get("shopId") ?? "")),
    ).subscribe({
      next: p =>  this._printers.set(p),
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }

  showCreateDialog(){
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }
}
