import { Component, input,  InputSignal, output, OutputEmitterRef } from '@angular/core';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { ShopAdminApiService } from '../../services/shop-admin-api/shop-admin-api.service';
import { FormsModule } from '@angular/forms';
import { Printer } from '../../domain/Printer';
import { NotificationService } from '../../services/notification/notification.service';
import { KlcsConfig } from '../../config/KlcsConfig';
import { ErrorResponse } from '../../domain/ErrorResponse';

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
    private notify: NotificationService,
  ){}

  createArticle() {
    this.articleDetails.Price = Math.floor(this.priceUi * 100);
    this.articleDetails.Printer = this.printer.Id === "" ? null : this.printer;
    const sub = this.shopAdminApi.createArticle(this.shopId(), this.articleDetails).subscribe({
      next: a => this.articleCreated.emit(a),
      error: (err: ErrorResponse) => this.notify.show({type: "error", duration: KlcsConfig.durationError, message: err.error.message}),
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
