# KRRU | Lan Catering System (KLCS)
This repository contains the implementation of KRRU's Catering System for Lan-Partys. This project aims to enable more than just KRRU Events organizing their Catering. Thus it can also be used for any other Catering use case.

**Requirements:**
* Database: Postgres
* Identity-Management: Keycloak

## Configurations
Configuration is done via environmental variables. The following variables can be used to configure KLCS.

### Environment variables
| Variable | Default | Description |
|----------|---------|-------------|
| LOGLVL | INFO | DEBUG, INFO, WARN, ERROR |
| API | :80 | Interface of API (default is port 80) |
| CORS_ORIGINS | | List of allowed Origins (CORS) - CSV Syntax |
| CORS_HEADERS | Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control | Allowed Headers |
| CORS_METHODS | POST,OPTIONS,GET,PUT,DELETE,PATCH | Allowed Methods |
| OIDC_JWKSURL | | Jwks URL of identity provider (Keycloak) |
| OIDC_AUTHORITY | | Oidc Authority |
| OIDC_CLIENTID | | Oidc ClientId |
| OIDC_ROLES_SYSADMIN | ADMIN | Role for Sysadmin |
| OIDC_ROLES_ACCOUNTMANAGER | ACCOUNT_MANAGER | Role for AccountManager |
| OIDC_ROLES_NOPRINT | NO_PRINT | User having this role will print no orders |
| DB_HOST | localhost | Database host (Postgres) |
| DB_PORT | 5432 | Database port |
| DB_USER | | Database user |
| DB_PASSWORD | | Database password |
| DB_NAME | postgres | Database name |
| DB_SSLMODE | verify-full | Database sslmode (disable, require, verify-ca, verify-full)|


## POS-Printer Support
KLCS offers a POS printer integration in addition to the purely digital system. Each shop in KLCS can have any number of printers, which can subsequently be assigned to articles.

When an article is sold, the order is printed on the configured printer. Should there be special cases where printing is only desired for certain sellers, this can be configured using an OIDC role.

The printer api can be used using the `printer-agent`. The Agent is a cli tool, that acts as a client for the printer api.

### Usage
The print-agent supports usb, as well as network connections to the corresponding printer. To be able to print orders, it is necessary, that the print-agent is able to estabish a network connection to klcs-backend (printer-api).

```text
/Projects/klcs/backend/print-agent$ go run main.go --help
Usage of .cache/go-build/c3/c3049deb6a54b98fe33184ca8aedb91ef4879375aadb7414030bfd1dd841d475-d/main:
  -klcshost string
    
  -loglvl string
         (default "INFO")
  -printer-id string
    
  -printer-netaddr string
    
  -printer-usbaddr string
    
  -timezone string
         (default "Europe/Vienna")
```
