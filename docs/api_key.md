# API Key

To use `mstodo`, you need to obtain a Microsoft API key for your Microsoft account.

## Steps

1. Sign in to the [Azure portal](https://portal.azure.com/).
2. Open the **"Azure Active Directory"** service (you may need to use the search bar at the top).
3. In the left-hand navigation pane, select **"App Registrations"** (it may be under the **"Manage"** heading).
4. Click the **"New registration"** button, near the top-left.
5. Populate the fields:
   - Enter `MS To Do CLI` into the **Name** field.
   - Select `Accounts in any organizational directory (Any Azure AD directory - Multitenant) and personal Microsoft accounts (e.g. Skype, Xbox)`.
   - Enter `http://localhost/oauth/callback` into the **Redirect URI** field.
6. Click **"Register"**
7. You will be redirected to the application's page.
8. Copy the **"Application (client) ID"**, and paste it into `.mstodo/config.yaml` after `client-id`. For example:

   ```yaml
   client-id: the-copied-application-client-id
   ```

9. In the left-hand navigation pane, select **"Certificates & secrets"** (it may be under the **"Manage"** heading).
10. Click **"New client secret"** under the **"Client secrets"** heading.
11. Add a description - e.g., "cli secret".
12. Copy the **"Value"** _(not the Secret ID)_, and paste it into `.mstodo/config.yaml` after `client-secret`. For example:

    ```yaml
    client-secret: the-copied-value
    ```

As a result, `~/.mstodo/config.yaml` should look like:

```yaml
client-id: the-copied-application-client-id
client-secret: the-copied-value
port: 12345 # The port you want mstodo to run on
auth-timeout: 120 # How long you want to wait until authentication times out
```
