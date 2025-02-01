import { Component, EventEmitter, Input, OnInit, Output, signal, Signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap, subscribeOn } from 'rxjs';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { CreateArticleDialogComponent } from "../../dialogs/create-article-dialog/create-article-dialog.component";
import { Printer } from '../../domain/Printer';
import { UpdateArticleDialogComponent } from "../../dialogs/update-article-dialog/update-article-dialog.component";

@Component({
  selector: 'klcs-shop-articles',
  imports: [
    CommonModule,
    CreateArticleDialogComponent,
    UpdateArticleDialogComponent
],
  templateUrl: './shop-articles.component.html',
  styleUrl: './shop-articles.component.css'
})
export class ShopArticlesComponent {
  @Input() categories: Signal<Map<string, Article[]>> = signal(new Map());
  @Output() articlesChanged: EventEmitter<void> = new EventEmitter();

  protected readonly CREATE_DIALOG_ID = "create-article-dialog"
  protected readonly EDIT_DIALOG_ID = "edit-article-dialog"

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private route: ActivatedRoute,
  ){}

  _printers: WritableSignal<Printer[]> = signal([])
  _articleDetails: WritableSignal<ArticleDetails> = signal(new ArticleDetails());

  deleteArticle(articleId: string) {
    if(confirm(`Do you realy want to delete Article?`)){
      const sub = this.shopAdminApi.deleteArticle(articleId).subscribe({
        next: _ => this.articlesChanged.emit(),
        error: err => console.error(err),
        complete: () => sub.unsubscribe(),
      });
    }
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
    this.refreshPrinters()
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }

  showUpdateDialog(articleId: string){
    this.refreshPrinters()
    const sub = this.shopAdminApi.getArticle(articleId).subscribe({
      next: artilce => {
        this._articleDetails.set(artilce)
        const dialog = document.getElementById(this.EDIT_DIALOG_ID) as HTMLDialogElement
        dialog.showModal();
      },
      error: err => console.error(err),
      complete: () => sub.unsubscribe(),
    })
  }
}
