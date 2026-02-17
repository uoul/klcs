export interface AppSettings {
  Version: string,
  Oidc: {
    JwksUrl: string;
    Authority: string;
    ClientId: string;
    Roles: {
      SysAdmin: string;
      AccountManager: string;
    };
  };
}
