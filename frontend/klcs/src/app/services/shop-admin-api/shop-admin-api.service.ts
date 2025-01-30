import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Shop } from '../../domain/Shop';
import { KlcsConfig } from '../../config/KlcsConfig';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ShopAdminApiService {

  constructor(
    private http: HttpClient,
  ) { }
}
