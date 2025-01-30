import { AuthConfig } from 'angular-oauth2-oidc';

export const KlcsConfig: { OAuth: AuthConfig, BackendRoot: string, RoleAdmin: string, RoleAccountManager: string} = {
  OAuth: {
    issuer: 'http://localhost:8081/realms/klcs',
    tokenEndpoint: 'http://localhost:8081/realms/klcs/protocol/openid-connect/token',
    redirectUri: window.location.origin + "/home",
    postLogoutRedirectUri: window.location.origin,
    clientId: 'klcs',
    responseType: 'code',
    scope: 'profile email',
    showDebugInformation: true,
  },
  BackendRoot: "http://localhost:8082",
  RoleAccountManager: "KLCS_ACCOUNT_MANAGER",
  RoleAdmin: "KLCS_ADMIN"
}
