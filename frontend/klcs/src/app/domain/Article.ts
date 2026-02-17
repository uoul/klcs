export class Article {
  constructor(
    public Id: string = "",
    public Name: string = "",
    public Description: string = "",
    public Price: number = 0,
    public StockAmount: number|null = null
  ){}
}
