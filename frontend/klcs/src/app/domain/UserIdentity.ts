export class UserIdentity {
  constructor(
    public username: string = "",
    public fullname: string = "",
    public firstname: string = "",
    public lastname: string = "",
    public email: string = "",
    public roles: string[] = [],
  ){}

  public static of(accessToken: string): UserIdentity {
    const userInfo = JSON.parse(atob(accessToken.split(".")[1]));
    return new UserIdentity(
      userInfo["preferred_username"] ?? "",
      userInfo["name"] ?? "",
      userInfo["given_name"] ?? "",
      userInfo["family_name"] ?? "",
      userInfo["email"] ?? "",
      userInfo["realm_access"]["roles"] ?? [],
    );
  }
}
