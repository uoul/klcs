import { Component, signal, WritableSignal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { firstValueFrom } from 'rxjs';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { CreateArticleDialogComponent } from "../../dialogs/create-article-dialog/create-article-dialog.component";
import { Printer } from '../../domain/Printer';
import { UpdateArticleDialogComponent } from "../../dialogs/update-article-dialog/update-article-dialog.component";
import { NotificationService } from '../../services/notification/notification.service';
import { SellerApiService } from '../../services/seller-api/seller-api.service';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'klcs-shop-articles',
  imports: [
    CommonModule,
    CreateArticleDialogComponent,
    UpdateArticleDialogComponent,
    TranslatePipe,
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
    protected translate: TranslateService,
  ){}

  _printers: WritableSignal<Printer[]> = signal([])
  _articleDetails: WritableSignal<ArticleDetails> = signal(new ArticleDetails());

  async deleteArticle(articleId: string) {
    if(confirm(this.translate.instant("components.shop-articles.DeletePrompt"))){
      try {
        await firstValueFrom(this.shopAdminApi.deleteArticle(articleId))
        await this.sellerApi.refreshShopDetails()
      } catch {}
    }
  }

  async refreshPrinters(shopId: string) {
    try {
      const printers = await firstValueFrom(this.shopAdminApi.getPrinters(shopId))
      this._printers.set(printers)
    } catch {}
  }

  showCreateDialog(){
    this.refreshPrinters(this.sellerApi.getShopDetails().Id)
    const dialog = document.getElementById(this.CREATE_DIALOG_ID) as HTMLDialogElement
    dialog.showModal();
  }

  async showUpdateDialog(articleId: string){
    this.refreshPrinters(this.sellerApi.getShopDetails().Id)
    try {
      const article = await firstValueFrom(this.shopAdminApi.getArticle(articleId))
      this._articleDetails.set(article)
      const dialog = document.getElementById(this.EDIT_DIALOG_ID) as HTMLDialogElement
      dialog.showModal();
    } catch {}
  }
}
