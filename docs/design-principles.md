# Dashboard Design Principles

## Vision

Build a dashboard that feels like a modern SaaS product rather than a traditional server administration panel.

The experience should resemble:

- Cloudflare
- Vercel
- Railway
- Linear
- GitHub
- Clerk

Users should immediately feel:

- Fast
- Clean
- Professional
- Trustworthy
- Minimal
- Powerful

---

# Core Design Philosophy

## Simple First

Avoid visual noise.

Every component must answer:

> Does this help the user make a decision?

If not, remove it.

---

## Information Density Without Clutter

Show important information.

Hide secondary information behind:

- Drawers
- Expandable Rows
- Dialogs
- Detail Panels

Example:

Do not show ASN, ISP, Country, User Agent directly in table.

Show:

```text
IP Address
Ban Time
Remaining Time
```

Click row → Open Details Panel.

---

## Modern SaaS Look

Avoid:

❌ Bootstrap Admin Panels

❌ Traditional Linux Control Panels

❌ Sharp Borders Everywhere

❌ Multiple Bright Colors

❌ Heavy Gradients

Use:

✅ Soft Shadows

✅ Rounded Corners

✅ Large White Space

✅ Clean Typography

✅ Subtle Animations

---

# Layout Structure

## Desktop Layout

```text
┌─────────────────────────────────────┐
│ Header                              │
├────────────┬────────────────────────┤
│ Sidebar    │ Main Content           │
│            │                        │
│            │                        │
└────────────┴────────────────────────┘
```

---

# Sidebar

Width:

```text
240px
```

Collapsed:

```text
72px
```

---

## Sidebar Menu

```text
Dashboard

Active Bans

Ban History

Live Requests

Analytics

Whitelist

Audit Logs

Settings
```

---

# Header

Height:

```text
64px
```

Contains:

- Search
- Theme Switch
- Notifications
- User Menu

---

# Content Width

Maximum:

```text
1600px
```

Centered layout.

Never stretch content edge-to-edge.

---

# Color Philosophy

## Light Mode

Background:

```text
#FFFFFF
```

Secondary Background:

```text
#F8FAFC
```

Card Background:

```text
#FFFFFF
```

Border:

```text
#E2E8F0
```

Primary Text:

```text
#0F172A
```

Secondary Text:

```text
#64748B
```

---

## Dark Mode

Background:

```text
#09090B
```

Secondary Background:

```text
#18181B
```

Card Background:

```text
#111827
```

Border:

```text
#27272A
```

Primary Text:

```text
#FAFAFA
```

Secondary Text:

```text
#A1A1AA
```

---

# Accent Color

Use only one primary accent.

Recommended:

```text
#3B82F6
```

Blue

Reasons:

- Security dashboards use blue well.
- Professional.
- Good contrast.
- Works in light and dark mode.

---

# Status Colors

## Success

```text
#22C55E
```

Green

Examples:

- Healthy
- Active
- Connected

---

## Warning

```text
#F59E0B
```

Amber

Examples:

- Elevated Traffic
- High Requests

---

## Danger

```text
#EF4444
```

Red

Examples:

- Banned
- Attack Detected
- Service Down

---

## Information

```text
#3B82F6
```

Blue

Examples:

- Statistics
- General Status

---

# Cards

## Style

Background:

```text
White / Dark Surface
```

Border Radius:

```text
16px
```

Padding:

```text
24px
```

Border:

```text
1px subtle border
```

Shadow:

```text
Very light shadow
```

Avoid:

```text
Heavy shadows
```

---

# Dashboard Cards

Example:

```text
Active Bans

145
+12 today
```

Layout:

```text
Metric
Label
Trend
```

Keep simple.

---

# Tables

Tables are the most important UI component.

---

## Table Design

Row Height:

```text
56px
```

Header:

```text
Sticky
```

Hover:

```text
Subtle background change
```

Selection:

```text
Accent border
```

---

## Active Bans Table

Show only:

```text
IP
Country
Ban Time
Remaining Time
Reason
Actions
```

Hide everything else.

---

# Drawers

Use right-side drawer.

Never open separate pages for details.

---

## Example

Click:

```text
192.168.1.1
```

Drawer opens:

```text
IP Details

Country
ASN
ISP

Request Count

429 Count

Ban History

Actions
```

---

# Charts

Keep charts minimal.

Avoid:

❌ 3D Charts

❌ Pie Charts Everywhere

❌ Excessive Colors

---

## Recommended

Line Charts

Area Charts

Bar Charts

---

# Animations

Purpose:

Guide attention.

Not decoration.

---

## Duration

Recommended:

```text
150ms - 250ms
```

---

## Use For

- Drawer Open
- Modal Open
- Notifications
- Table Updates
- Theme Changes

---

## Avoid

- Bouncing Elements
- Continuous Animations
- Floating Widgets

---

# Theme Switching

Support:

- Light
- Dark
- System

Save preference.

---

# Dark Mode Priority

Design dark mode first.

Reason:

Security dashboards are primarily used by engineers.

Most engineers use dark mode.

---

# Typography

Font:

## Inter

Recommended.

Alternative:

## Geist

---

# Font Sizes

Page Title:

```text
30px
```

Section Title:

```text
20px
```

Card Metric:

```text
32px
```

Body:

```text
14px
```

Small Text:

```text
12px
```

---

# Empty States

Never show:

```text
No Data
```

Instead:

```text
No active bans detected.
Traffic appears normal.
```

---

# Loading States

Use:

- Skeleton Loaders
- Animated Placeholders

Avoid:

```text
Loading...
```

---

# Notifications

Position:

```text
Top Right
```

Use Toasts.

Examples:

```text
IP successfully unbanned
```

```text
Whitelist updated
```

```text
Configuration saved
```

---

# Responsive Design

## Desktop

Primary experience.

---

## Tablet

Fully supported.

---

## Mobile

Read-only friendly.

Allow:

- View bans
- Search
- Unban

Do not optimize for heavy administration.

---

# Inspiration Sources

Use design references from:

- Cloudflare Dashboard
- Vercel Dashboard
- Railway Dashboard
- GitHub Settings
- Linear
- Clerk

Avoid design references from:

- cPanel
- Plesk
- Old Grafana Versions
- Legacy Admin Templates

---

# Final UI Goals

The dashboard should feel:

✅ Modern

✅ Fast

✅ Professional

✅ Minimal

✅ Trustworthy

✅ Premium

✅ Easy To Use

✅ Beautiful In Both Light And Dark Mode

Users should feel like they are using a modern SaaS platform, not a server management tool.

---

# Reference Images

Reference Images store inside Refrance-site-UI-images folder