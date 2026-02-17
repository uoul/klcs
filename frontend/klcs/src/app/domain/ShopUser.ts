import {Role} from "./Role";

export class ShopUser {
  constructor(
    public Id: string = "",
    public Name: string = "",
    public Username: string = "",
    public ShopRoles: Role[] = [],
  ) {
  }
}
