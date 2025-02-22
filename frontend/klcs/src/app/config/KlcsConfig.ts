export const KlcsConfig = {
  OAuth: {
    issuer: 'https://kc.uoul.net/realms/klcs',
    tokenEndpoint: 'https://kc.uoul.net/realms/klcs/protocol/openid-connect/token',
    redirectUri: window.location.origin,
    postLogoutRedirectUri: window.location.origin,
    clientId: 'klcs',
    responseType: 'code',
    scope: 'profile email',
    showDebugInformation: true,
  },
  BackendRoot: "http://localhost:8082",
  
  KlcsRoleAccountManager: "KLCS_ACCOUNT_MANAGER",
  KlcsRoleAdmin: "KLCS_ADMIN",

  ShopRoleAdmin: "ADMIN",

  durationShort: 5000,
  durationMedium: 10000,
  durationLong: 15000,
}
