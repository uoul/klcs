export class HistoryArticle {
  constructor(
    public Id: string = "",
    public Name: string = "",
    public Description: string = "",
    public Pieces: number = 0,
    public PrinterAck: boolean = false,
  ){}
}
