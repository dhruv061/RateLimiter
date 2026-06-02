# new-features-1.md

# Additional Feature Requirements

This document contains additional requirements discovered after the initial feature specification.

These requirements must be considered mandatory and implemented alongside the features defined in `features.md`.

---

# 1. Global Time & Date Filtering

## Objective

Allow users to view data for a specific time period across the entire dashboard.

---

## Global Filter Bar

Add a global filter section in the dashboard header.

Visible on all pages.

Example:

```text
Domain ▼

Today ▼

Refresh
```

---

## Supported Time Ranges

* Last 15 Minutes
* Last 30 Minutes
* Last 1 Hour
* Last 6 Hours
* Last 12 Hours
* Last 24 Hours
* Last 7 Days
* Last 30 Days
* Last 90 Days
* Custom Range

---

## Custom Range

Allow user to select:

```text
Start Date
Start Time

End Date
End Time
```

Example:

```text
2026-06-01 00:00

to

2026-06-02 23:59
```

---

## Affected Modules

Time filtering must affect:

* Dashboard Metrics
* Active Bans
* Ban History
* Analytics
* Traffic Charts
* Top Offenders
* Audit Logs
* Live Request Logs

---

# 2. Multi-Domain Support

## Objective

Dashboard should support multiple domains.

A user may manage:

```text
example.com

api.example.com

test.example.com

domain.com
```

from a single dashboard.

---

# Domain Management

## Remove Environment Variable Dependency

Do not hardcode domains using environment variables.

Instead:

Create domain management within the dashboard.

---

# Domains Page

New menu item:

```text
Domains
```

---

## Add Domain Wizard

User can add a domain using a setup wizard.

---

### Required Information

#### Basic Information

```text
Domain Name
```

Example:

```text
example.com
```

---

#### Nginx Information

```text
Access Log Path

Error Log Path

Blocked IP File Path
```

Examples:

```text
/var/log/nginx/example_access.log

/var/log/nginx/example_error.log

/etc/nginx/example_blocked.conf
```

---

#### Fail2Ban Information

```text
Fail2Ban Jail Name
```

Examples:

```text
nginx-429

nginx-limit-req
```

---

#### Server Information

```text
Server Name
Description
```

Optional.

---

# Domain Validation

Before saving:

Validate:

* Access log exists
* Error log exists
* Block file exists
* Fail2Ban jail exists

Show validation results.

---

# Domain Selection

## Global Domain Selector

Add dropdown in top navigation.

Example:

```text
┌────────────────────────┐
│ Domain: example.com ▼  │
└────────────────────────┘
```

---

# Dashboard Behaviour

All data displayed should be filtered based on the selected domain.

Examples:

### Domain Selected

```text
example.com
```

Dashboard shows:

* example.com bans
* example.com logs
* example.com analytics

Only.

---

### Domain Selected

```text
api.example.com
```

Dashboard shows:

* api.example.com bans
* api.example.com logs
* api.example.com analytics

Only.

---

# Remember Domain

Persist selected domain.

Store:

```text
Local Storage
```

or

```text
User Preferences
```

---

# 3. Component Help System

## Objective

Users should understand every widget and metric.

---

# Help Icon

Every major component must include:

```text
ⓘ
```

help icon.

---

# Location

Top-right corner of component header.

Example:

```text
Active Bans                         ⓘ
```

---

# Behaviour

Click icon:

Show tooltip or modal.

Explain:

* What component does
* What metrics mean
* How data is calculated
* Why user should care

---

# Examples

## Active Bans

Description:

```text
Displays all currently banned IP addresses.

These IPs are currently blocked from accessing the selected domain.

Use this table to review, investigate, or manually unban IPs.
```

---

## 429 Responses

Description:

```text
Shows how many requests were rejected due to rate limiting.

A higher value may indicate abusive traffic or attack attempts.
```

---

## Top Offenders

Description:

```text
Shows IP addresses generating the highest number of requests or rate limit violations.

Useful for identifying abusive clients.
```

---

# Help Content Requirements

Every component must have:

* Title
* Description
* Metric Explanation
* Usage Guidance

---

# 4. Enhanced Dashboard Context

## Dashboard Header

Show:

```text
Selected Domain

Selected Time Range

Last Refresh Time
```

Example:

```text
example.com

Last 24 Hours

Updated 5 seconds ago
```

---

# 5. Domain Isolation

## Critical Requirement

All queries must be scoped to the selected domain.

No component should ever mix data from multiple domains.

---

# 6. Future Expansion Support

Database design should support:

```text
1 User
Many Domains

1 Domain
Many Logs

1 Domain
Many Bans

1 Domain
Many Analytics Records
```

---

# Required New Pages

Add:

```text
Domains
```

to sidebar navigation.

Updated menu:

```text
Dashboard

Domains

Active Bans

Ban History

Live Requests

Analytics

Whitelist

Audit Logs

Settings
```

---

# Success Criteria

A user should be able to:

* Add a new domain
* Validate domain configuration
* Select a domain from the header
* Filter all data by domain
* Filter all data by date/time range
* Understand every widget through help icons
* Switch domains without page reload
* View isolated analytics for each domain

These features are mandatory and should be treated as part of the core product requirements.
