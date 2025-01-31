import {Article} from "./Article";

export class ShopDetails {
    constructor(
        public Id: string = "",
        public Name: string = "",
        public UserRoles: string[] = [],
        public Categories: Map<string, Article[]> = new Map<string, Article[]>()
    ){}
}
