---
type: reality
status: active
owner: slashben
scope: repo
related_code:
  - notifications/runtimeincidents.go
  - notifications/runtimeincidents_test.go
---

# Runtime Incident Actor

Lets runtime incident notifications (Slack, Teams, Jira, webhooks and the other
alert channels) show who triggered an incident, for example the Kubernetes user
who ran "Exec to pod". This library change defines the notification contract
only. Populating it lives in event-ingester-service and rendering it in
users-notification-service, both in armosec/backend.

## The problem

Admission alerts already record the requesting Kubernetes user.
`armotypes.AdmissionAlert.UserInfo` holds the username and groups, and the
backend's stored runtime incident embeds `armotypes.RuntimeAlert`, so it has
the user too. The notification payload had no field for it, so the user was
dropped before any alert channel could show it.

## What this adds

Two fields on `NewRuntimeIncident`:

```go
UserName   string   `json:"userName,omitempty"`
UserGroups []string `json:"userGroups,omitempty"`
```

Nothing existing changed. No field was renamed or removed.

## Why the JSON tags must match kubescape/messaging

event-ingester-service sends `RuntimeIncidentIngesterOnFinishedMessage` from
`github.com/kubescape/messaging/pulsar/common/kdr`, which gained the same two
fields with the same tags in v0.0.114. users-notification-service does not copy
that message field by field into `NewRuntimeIncident`. It passes the whole
message through as the notification data, which is decoded as JSON. The tags are
therefore the contract: a renamed tag on either side silently drops the actor.
`TestNewRuntimeIncidentActorJSON` decodes a payload in the messaging shape and
pins the tags.

## Why omitempty

Only admission alerts have a Kubernetes user. Node-agent rules (process, file,
network) and CDR alerts leave both fields empty, and `omitempty` keeps those
payloads byte-identical to before. `TestNewRuntimeIncidentActorOmittedWhenEmpty`
pins this. Templates should render the actor row only when `UserName` is set.

## No bson tags

`AGENTS.md` asks for `json` and `bson` tags together. `NewRuntimeIncident` is a
notification payload that is never stored in MongoDB, and none of its existing
fields carry `bson` tags, so the new fields follow the struct and use `json`
only.

## What this does not cover

Node-agent and CDR alerts have no Kubernetes user. An actor for those (process
or UID on the node, cloud identity for CDR) would be a separate change.
