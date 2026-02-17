import {Printer} from "./Printer";

export class ArticleDetails {
  constructor(
    public Id: string = "",
    public Name: string = "",
    public Description: string = "",
    public Price: number = 0,
    public StockAmount: number | null = null,
    public Category: string = "",
    public Printer: Printer | null = null,
  ) {
  }
}
