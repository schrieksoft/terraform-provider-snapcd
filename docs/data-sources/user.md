---
page_title: "snapcd_user Data Source - snapcd"
subcategory: "Identity Access Management"
description: |-
  Use this data source to access information about an existing User in Snap CD.
---

# snapcd_user (Data Source)

Use this data source to access information about an existing User in Snap CD.


## Example Usage

```terraform
data "snapcd_user" "myuser" {
  name = "somebody@somewhere.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `user_name` (String) Unique name of the user.

### Read-Only

- `access_failed_count` (Number) The number of failed access attempts.
- `concurrency_stamp` (String) Used to handle concurrency checks.
- `email` (String) User's email address.
- `email_confirmed` (Boolean) Whether the user's email has been confirmed.
- `id` (String) Unique ID of the User.
- `is_active` (Boolean) Indicates whether the user is active.
- `lockout_enabled` (Boolean) Indicates whether lockout is enabled for the user.
- `lockout_end` (String) The date and time when the lockout ends (if any).
- `normalized_email` (String) Normalized email address used for consistency.
- `normalized_user_name` (String) Normalized user name used for consistency.
- `password_hash` (String, Sensitive) Hashed password of the user.
- `phone_number` (String) Phone number of the user.
- `phone_number_confirmed` (Boolean) Whether the phone number is confirmed.
- `security_stamp` (String) Security stamp used to identify changes to the user's security info.
- `two_factor_enabled` (Boolean) Indicates if two-factor authentication is enabled.
