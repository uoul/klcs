import { Component, EventEmitter, input, Input, InputSignal, OnInit, output, Output, OutputEmitterRef, signal, Signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';
import { CommonModule } from '@angular/common';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { ActivatedRoute } from '@angular/router';
import { mergeMap, subscribeOn, take } from 'rxjs';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { CreateArticleDialogComponent } from "../../dialogs/create-article-dialog/create-article-dialog.component";
import { Printer } from '../../domain/Printer';
import { UpdateArticleDialogComponent } from "../../dialogs/update-article-dialog/update-article-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { ShopDetails } from '../../domain/ShopDetails';

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

  protected readonly CREATE_DIALOG_ID = "create-article-dialog"
  protected readonly EDIT_DIALOG_ID = "edit-article-dialog"

  constructor(
    private shopAdminApi: ShopAdminApiService,
    protected sellerApi: SellerApiService,
    private notify: NotificationService,
  ){}

  _printers: WritableSignal<Printer[]> = signal([])
  _articleDetails: WritableSignal<ArticleDetails> = signal(new ArticleDetails());

  deleteArticle(articleId: string) {
    if(confirm(`Do you realy want to delete Article?`)){
      const sub = this.shopAdminApi.deleteArticle(articleId).subscribe({
        next: _ => this.sellerApi.refreshShopDetails(),
        error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
        complete: () => sub.unsubscribe(),
      });
    }
  }

  refreshPrinters(shopId: string) {
    const sub = this.shopAdminApi.getPrinters(shopId).subscribe({
      next: p =>  this._printers.set(p),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }

  showCreateDialog(){
    this.refreshPrinters(this.sellerApi.getShopDetails().Id)
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }

  showUpdateDialog(articleId: string){
    this.refreshPrinters(this.sellerApi.getShopDetails().Id)
    const sub = this.shopAdminApi.getArticle(articleId).subscribe({
      next: artilce => {
        this._articleDetails.set(artilce)
        const dialog = document.getElementById(this.EDIT_DIALOG_ID) as HTMLDialogElement
        dialog.showModal();
      },
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
      complete: () => sub.unsubscribe(),
    })
  }
}
