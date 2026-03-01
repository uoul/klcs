import { Component, input,  InputSignal, output, OutputEmitterRef, signal, WritableSignal } from '@angular/core';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Printer } from '../../domain/Printer';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'klcs-create-article-dialog',
  imports: [
    FormsModule,
    TranslatePipe,
  ],
  templateUrl: './create-article-dialog.component.html',
  styleUrl: './create-article-dialog.component.css'
})
export class CreateArticleDialogComponent {
  nonePrinter: Printer = { Id: "", Name: "-", Connected: false };

  dialogId: InputSignal<string> = input.required<string>();
  shopId: InputSignal<string> = input.required<string>();
  printers: InputSignal<Printer[]> = input([this.nonePrinter]);

  dialogClosed: OutputEmitterRef<void> = output();
  articleCreated: OutputEmitterRef<ArticleDetails> = output();

  articleDetails: ArticleDetails = new ArticleDetails()
  priceUi: number = 0.0;
  printer: Printer = this.nonePrinter;
  isActive: WritableSignal<boolean> = signal(false)

  constructor(
    private shopAdminApi: ShopAdminApiService,
    private notify: NotificationService,
    protected translate: TranslateService,
  ){}

  async createArticle() {
    if(!this.isActive()){
      this.isActive.set(true)
      this.articleDetails.Price = Math.floor(this.priceUi * 100);
      this.articleDetails.Printer = this.printer.Id === "" ? null : this.printer;
      try {
        const createdArticle = await firstValueFrom(this.shopAdminApi.createArticle(this.shopId(), this.articleDetails))
        this.articleCreated.emit(createdArticle)
      } finally {
        this.init()
        this.isActive.set(false)
      }
    }    
  }

  init(){
    this.articleDetails = new ArticleDetails()
    this.priceUi = 0.0
    this.printer = this.nonePrinter;
  }
}
