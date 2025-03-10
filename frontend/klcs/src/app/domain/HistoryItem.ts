import { HistoryArticle } from "./HistoryArticle";

export class HistoryItem {
  constructor(
    public TransactionId: string = "",
    public Timestamp: Date = new Date(),
    public Description?: string,
    public AccountHolder?: string,
    public Articles: HistoryArticle[] = [],
  ){}
}