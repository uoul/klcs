import { AuthConfig } from 'angular-oauth2-oidc';

export const authCodeFlowConfig: AuthConfig = {
  issuer: 'http://localhost:8081/realms/klcs',
  tokenEndpoint: 'http://localhost:8081/realms/klcs/protocol/openid-connect/token',
  redirectUri: window.location.origin + "/home",
  postLogoutRedirectUri: window.location.origin,
  clientId: 'klcs',
  responseType: 'code',
  scope: 'profile email',
  showDebugInformation: true,
};