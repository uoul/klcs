export class Order {
  constructor(
      public Type: string = "",
      public Description: string = "",
      public AccountId: string | null = null,
      public Sum: number | undefined = undefined,
      public Articles: {[name: string]: number} = {},
  ) {}
}
