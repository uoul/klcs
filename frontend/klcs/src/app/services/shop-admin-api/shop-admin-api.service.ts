import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {catchError, Observable, of} from "rxjs";
import { Article } from '../../domain/Article';
import { ArticleDetails } from '../../domain/ArticleDetails';
import { Printer } from '../../domain/Printer';
import { Role } from '../../domain/Role';
import { ShopUser } from '../../domain/ShopUser';
import { User } from '../../domain/User';
import { KlcsConfig } from '../../config/KlcsConfig';

@Injectable({
  providedIn: 'root'
})
export class ShopAdminApiService {

  constructor(
    private http: HttpClient,
  ) { }

  public getArticlesForShop(shopId: string): Observable<Article[]> {
    return this.http.get<Article[]>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/articles`);
  }

  public createArticle(shopId: string, article: ArticleDetails): Observable<ArticleDetails> {
    return this.http.post<ArticleDetails>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/articles`, article);
  }

  public getArticle(articleId: string): Observable<ArticleDetails> {
    return this.http.get<ArticleDetails>(`${KlcsConfig.BackendRoot}/api/v1/articles/${articleId}`)
  }

  public updateArticle(article: ArticleDetails): Observable<ArticleDetails> {
    return this.http.patch<ArticleDetails>(`${KlcsConfig.BackendRoot}/api/v1/articles/${article.Id}`, article)
  }

  public deleteArticle(articleId: string): Observable<Object> {
    return this.http.delete(`${KlcsConfig.BackendRoot}/api/v1/articles/${articleId}`)
  }

  public getPrinters(shopId: string): Observable<Printer[]> {
    return this.http.get<Printer[]>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/printers`)
  }

  public createPrinterForShop(shopId: string, printer: Printer): Observable<Printer> {
    return this.http.post<Printer>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/printers`, printer)
  }

  public deletePrinter(printerId: string): Observable<Object> {
    return this.http.delete(`${KlcsConfig.BackendRoot}/api/v1/printers/${printerId}`)
  }

  public getUsers(): Observable<User[]> {
    return  this.http.get<User[]>(`${KlcsConfig.BackendRoot}/api/v1/users`)
  }

  public getRoles(): Observable<Role[]> {
    return this.http.get<Role[]>(`${KlcsConfig.BackendRoot}/api/v1/roles`)
  }

  public getUsersForShop(shopId: string): Observable<ShopUser[]> {
    return this.http.get<ShopUser[]>(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/users`)
  }

  public addUserRoleForShop(shopId: string, userId:string, role:Role): Observable<Object> {
    return this.http.post(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/users/${userId}/roles`, role)
  }

  public deleteUserRoleForShop(shopId: string, userId: string, roleId: string): Observable<Object> {
    return this.http.delete(`${KlcsConfig.BackendRoot}/api/v1/shops/${shopId}/users/${userId}/roles/${roleId}`)
  }
}
