# new-features-2.md

# Fail2Ban First Setup Experience

## Objective

The dashboard should not assume that domains are already configured.

The first experience should guide the user through:

1. Detecting Nginx domains.
2. Configuring Fail2Ban.
3. Verifying Nginx changes.
4. Activating protection.
5. Showing domain analytics.

The product should feel like a setup wizard rather than a technical administration panel.

Reference implementation should follow the setup process described in:

```text
rate-limit-setup.md
```

This document is the source of truth for automatic Fail2Ban setup generation.

---

# 1. Empty State Experience

## First Launch

When no configured domains exist:

DO NOT show:

* Empty dashboard
* Empty charts
* Empty tables

Instead show a setup screen.

---

## Screen

```text
Welcome to ShieldWatch

No protected domains configured.

Let's connect your first domain and enable protection.

[ Configure Domain ]
```

---

# 2. Configure Domain Flow

Click:

```text
Configure Domain
```

User is redirected to:

```text
Setup Wizard
```

---

# 3. Setup Wizard Steps

## Step 1

Detect Fail2Ban Status

---

### Show

```text
Fail2Ban Status

● Running
```

or

```text
Fail2Ban Status

● Not Running
```

---

### Information

Display:

* Service Status
* Version
* Active Jails
* Installed / Not Installed

---

### Actions

If not installed:

```text
Install Fail2Ban
```

If installed:

```text
Continue
```

---

# 4. Existing Configurations

Display current configured domains.

Example:

```text
example.com

Status: Active
```

```text
api.example.com

Status: Active
```

---

# 5. Add New Domain

Button:

```text
Add Domain
```

---

# 6. Domain Discovery

Automatically scan:

```text
/etc/nginx/sites-enabled
```

```text
/etc/nginx/conf.d
```

```text
/etc/nginx/sites-available
```

---

## Show List

Example:

```text
example.com

api.example.com

test.example.com
```

---

User selects domain.

---

# 7. Auto Configuration Generation

After selecting domain:

Generate all required Fail2Ban files automatically.

Based on:

```text
rate-limit-setup.md
```

Generate:

* Filter
* Action
* Jail
* Block File
* Domain Config

Using domain-specific names.

---

## Example

Domain:

```text
example.com
```

Generate:

```text
example-com-limit-req
```

```text
example-com-429
```

---

## Block File

Generate:

```text
/etc/nginx/example-com_blocked.conf
```

instead of:

```text
/etc/nginx/fail2ban_blocked.conf
```

---

## Jail Names

Generate:

```text
example-com-limit-req
```

```text
example-com-429
```

instead of shared names.

---

# 8. Automatic Shell Script Generation

Dashboard generates a shell script.

Example:

```bash
setup-example-com.sh
```

Based on:

```text
rate-limit-setup.md
```

The user should not manually create files.

---

# 9. Nginx Configuration Step

After generating configuration:

Show:

```text
Step 4 of 5

Update Nginx Configuration
```

---

Display generated Nginx changes.

Beautiful code viewer.

Copy button.

Download button.

---

## UI Example

```text
┌──────────────────────────────┐

Required Nginx Changes

[ Copy ]

-----------------------------

generated config here

-----------------------------

└──────────────────────────────┘
```

---

# 10. User Verification

After updating Nginx:

Button:

```text
I Have Updated Nginx
```

---

# 11. Validation

Dashboard automatically validates:

* Nginx config exists
* Required directives exist
* Block file exists
* Jail exists
* Fail2Ban active

---

# 12. Success Screen

If validation succeeds:

Show full-screen success experience.

---

## Animation

Use:

* Fireworks
* Confetti
* Glow animation

---

## Message

```text
Protection Enabled

example.com is now protected.

Rate limiting and automatic IP banning are active.
```

---

# 13. Domain Status

Domain becomes visible in:

```text
Domain Selector
```

---

## Status

```text
● Active
```

Green.

---

# 14. Domain Cards

Domains page should show:

### Domain

```text
example.com
```

### Protection Status

```text
Active
```

### Fail2Ban Status

```text
Running
```

### Jails

```text
2 Active
```

### Last Ban

```text
5 minutes ago
```

---

# 15. Dashboard Access Rules

If domain is not activated:

Hide:

* Analytics
* Active Bans
* Ban History

Show setup wizard instead.

---

# 16. Reconfigure Domain

Domain card actions:

```text
View Details

Edit

Reconfigure

Delete
```

---

# 17. Delete Domain

Removing domain should:

* Remove dashboard configuration
* Remove database records

Optional:

```text
Remove Fail2Ban files from server
```

with confirmation.

---

# Success Criteria

A new user should be able to:

1. Install Fail2Ban.
2. Detect Nginx domains.
3. Select a domain.
4. Auto-generate configuration.
5. Apply Nginx changes.
6. Verify setup.
7. Activate protection.
8. Immediately view analytics.

Without manually reading setup documentation.

The dashboard should function as a guided setup experience rather than a configuration management panel.

# 18. Safe Domain Removal Workflow

## Objective

Removing a domain should never immediately delete configuration.

The dashboard must guide the user through a safe removal process to avoid:

* Broken Nginx configuration
* Orphaned Fail2Ban jails
* Leftover block files
* Invalid server state

---

# Remove Domain Action

When user clicks:

```text
Delete Domain
```

Do NOT delete immediately.

Show confirmation wizard.

---

# Step 1 - Warning Screen

Display:

```text
Warning

You are about to remove protection for:

example.com

This action will:

• Remove dashboard configuration
• Remove Fail2Ban jails
• Remove generated filter files
• Remove generated action files
• Remove generated block files
• Stop protection for this domain

Before continuing you must remove the generated Nginx configuration from your server.
```

Actions:

```text
Cancel

Continue
```

---

# Step 2 - Nginx Cleanup Instructions

Display generated configuration previously added during setup.

Example:

```nginx
limit_req zone=example-com_rate_limit burst=5 nodelay;
limit_req_status 429;

if ($example_com_blocked) {
    return 403;
}
```

---

## UI

```text
Required Nginx Cleanup

Remove the following configuration from your Nginx domain configuration.

[ Copy ]

--------------------------------

Generated nginx configuration

--------------------------------
```

---

## User Confirmation

Checkbox:

```text
☐ I have removed the configuration from Nginx.
```

Button:

```text
Validate
```

---

# Step 3 - Validation

Dashboard validates:

### Nginx Config

Verify:

```text
Generated configuration no longer exists.
```

---

### Fail Validation

If configuration still exists:

```text
Configuration still detected.

Please remove the generated configuration and try again.
```

Prevent deletion.

---

# Step 4 - Remove Fail2Ban Configuration

After validation succeeds:

Display summary.

Files that will be removed:

```text
/etc/fail2ban/jail.d/example-com.local

/etc/fail2ban/filter.d/example-com-429.conf

/etc/fail2ban/action.d/example-com-block.conf

/etc/nginx/example-com_blocked.conf
```

---

# User Confirmation

```text
☐ Remove generated Fail2Ban files
☐ Remove generated block file
☐ Restart Fail2Ban automatically
☐ Reload Nginx automatically
```

Button:

```text
Delete Protection
```

---

# Step 5 - Automated Cleanup

Dashboard executes generated cleanup script.

Example:

```bash
cleanup-example-com.sh
```

Operations:

* Remove jail file
* Remove filter file
* Remove action file
* Remove block file
* Reload Nginx
* Restart Fail2Ban
* Validate cleanup

---

# Step 6 - Post Cleanup Validation

Verify:

### Fail2Ban

Domain jail no longer exists.

### Nginx

Generated block file removed.

### Dashboard

Domain configuration removed.

---

# Failure Handling

If cleanup partially fails:

Show detailed error.

Example:

```text
Cleanup Failed

The following resources still exist:

• example-com.local
• example-com_blocked.conf

Please resolve manually.
```

Do NOT remove the domain record.

---

# Step 7 - Final Confirmation

If cleanup succeeds:

Show success animation.

Use:

* Confetti
* Fireworks
* Success Glow

---

## Message

```text
Protection Removed

example.com has been successfully removed.

All generated Fail2Ban resources have been cleaned up.

The domain is no longer protected by this dashboard.
```

---

# Database Cleanup

Only after successful validation:

Remove:

* Domain record
* Active bans
* Domain analytics
* Domain settings
* Domain audit records

Optional:

Keep historical records for audit purposes.

---

# Safety Rules

Never allow direct deletion.

Always require:

1. Nginx cleanup.
2. Validation.
3. Fail2Ban cleanup.
4. Validation.
5. Final deletion.

This prevents broken configurations and accidental loss of protection.
