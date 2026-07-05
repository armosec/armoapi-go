---
type: feature
status: in-progress
owner: alonliwsky
scope: repo
related_code:
  - notifications/integrationhealth.go
  - notifications/collaborationconfig.go
---

# Integration Token Health — shared model

Part of the cross-service integration-token-health project: OAuth ticketing
integrations (Linear, Jira Cloud, GitHub) silently break when refresh tokens
expire while the platform keeps reporting "connected". This library change
defines the shared health contract; detection, cron, UI, and notifications
live in the consuming services.

## What this adds

1. **`Degraded` value** on `IntegrationConnectionStatus` (`connected | disconnected | degraded`).
   `degraded` = the integration is configured but its credentials are unusable;
   user action (OAuth reconnect) is required.
2. **Health attributes contract** on `CollaborationConfig.Attributes`
   (`notifications/integrationhealth.go`):

   | Attribute key | Meaning |
   |---|---|
   | `healthStatus` | `"degraded"`; **absence = healthy** |
   | `healthLastChecked` | RFC3339 time of last health evaluation |
   | `healthLastError` | reason for the last failed evaluation |
   | `healthDegradedSince` | RFC3339 time of first transition into degraded |

3. **Typed helpers** so consumers never touch raw keys: `GetHealth()`,
   `SetHealthDegraded(reason)` (idempotent — keeps the original
   `healthDegradedSince` on re-mark), `SetHealthChecked()`, `ClearHealth()`.

## Why Attributes instead of typed struct fields

`CollaborationConfig` is read-modify-written by several services
(cadashboardbe, config-service, users-notification-service). A service built
against an older armoapi-go silently drops unknown *typed* fields when it
deserializes and re-serializes the record — a typed degraded flag would
flicker back to healthy whenever any lagging writer rewrote the config.
Attributes-map entries survive those round-trips with zero version
coordination, and remain queryable in Mongo (`attributes.healthStatus`).

## Consumers

- **cadashboardbe** — writes health on `RefreshTokenError` (reactive) and from
  the daily health check (proactive); reads it in the integrations status +
  configV2 endpoints. Clears on OAuth reconnect.
- **users-notification-service** — reads the 424 `integration_token_expired`
  contract derived from this state (does not write health itself).
- **armo-ui** — renders `degraded` badge/banner from the status endpoints.
