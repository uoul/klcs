# KRRU | Lan Catering System (KLCS)
This repository contains the implementation of KRRU's Catering System for Lan-Partys. This project aims to enable more than just KRRU Events organizing their Catering. Thus it can also be used for any other Catering use case.

**KLCS consists of three parts that are needed:**
* Database: Postgres
* Identity-Management: Keycloak
* KLCS: Frontend and Backend (Docker-Image)

## Configurations
Configuration is done via two methods
1. Via Environmental variables
2. Via KlcsConfig.ts (in Frontend)

### Environment variables
| Variable | Default | Description |
|----------|---------|-------------|
| KLCS_LOG_LVL | INFO | OFF, TRACE, DEBUG, INFO, WARNING, ERROR, FATAL |
| KLCS_JWKS_URI | "" | JWKS-Url of KeyCloak, used for authorization |
| KLCS_DB_HOST | "localhost" | Database Host (Postgres) |
| KLCS_DB_PORT | 5432 | Database Port |
| KLCS_DB_USER | "" | Database user |
| KLCS_DB_PW | "" | Database user's password |
| KLCS_DB_NAME | postgres | Database name |
| KLCS_DB_SSL | verify-full | Database flag for ssl (disable, allow, prefer, require, verify-ca, verify-full) |
