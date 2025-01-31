export class AccountDetails {
  constructor(
    public Id: string = "",
    public HolderName: string = "",
    public Locked: boolean = false,
    public Balance: number = 0,
  ){}
}
