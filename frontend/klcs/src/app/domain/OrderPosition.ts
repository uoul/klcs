import {Article} from "./Article";

export class OrderPosition {
  constructor(
    public article: Article = new Article(),
    public count: number = 0,
  ) {}
}
