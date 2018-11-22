# addons-template-service

The only requirement to have a custom Bitrise addon server is that the server must handle a `/provision` endpoint with the methods: `POST`, `PUT`, `DELETE`. It is required to be able to authenticate a Bitrise app on the server. Also this is how you can export environment variables into the build. (For example the url of this server, and an access token for the app). You can have as many addon server logic as you need, however this provision endpoint must use different authentication than the others, it must be the shared token that you will give to Bitrise to add your addon to the addon service.

We will need these infos when you want to add a new addon: https://github.com/bitrise-team/bitrise-addon-service/blob/master/_template/sample-addon.yml

# Authentication

All 3 endpoint requires authentication. This must be the shared token from the shared configuration. The header key must be `Authentication` and the value can be the shared token itself. Beside of the 3 provision handler, all of your other endpoints must use different authentication.

Example authentication header key-value pair:
 - `"Authentication" : "bitrise-shared-token"`

---

# Provision of a new app

> The server creates a new record or updates an existing one with the appslug, to store provision state of the app. Also store a unique token for the appslug that will be used for the requests that are from a Bitrise build and calls this server. Also store the received plan, so you can have a service that can use specified parameters/limits by the plan. Finally sends back the list of environment variables that will be exported in all of the builds on Bitrise for the app.

**Method**: POST

**URL**: `/provision`

**Content**:
```
{
    "plan": "free",
    "app_slug": "app-slug-123",
    "api_token": "public-API-token"
}
```

### Success Response

- Status: 200 Success

- Content:
    ```
    {
        "envs": [
            {
                "key": "MYADDON_HOST_URL",
                "value": "https://my-addon.url"
            },
            {
                "key": "MYADDON_AUTH_SECRET",
                "value": "verysecret"
            }
        ]
    }
    ```

*At this point in the next build you can use the envs $MYADDON_HOST_URL and $MYADDON_AUTH_SECRET.*

*For example:*

*POST https://$MYADDON_HOST_URL/store-cache (Authentication: $MYADDON_AUTH_SECRET)*

### Error Response

- Status: 403 Unauthorized

- Status: 500 InternalServerError

  - Content:
  ```
  {
      "error": "error message..."
  }
  ```

---

# Send plan update to an app

> Overwrite the plan that you saved already with the one that is in this request. This way Bitrise can update your addon if there was a plan change for any reason.

**Method**: PUT

**URL**: `/provision/{app_slug}`

**Content**:
```
{
    "plan": "free"
}
```

### Success Response

- Status: 200 Success

### Error Response

- Status: 403 Unauthorized

- Status: 500 InternalServerError

  - Content:
  ```
  {
      "error": "error message..."
  }
  ```

---

# Delete provision of an app

> Delete the app's provisioned state, so the calls are pointed to this service will be rejected in the Bitrise build.

**Method**: DELETE

**URL**: `/provision/{app_slug}`

### Success Response

- Status: 200 Success

### Error Response

- Status: 403 Unauthorized

- Status: 500 InternalServerError

  - Content:
  ```
  {
      "error": "error message..."
  }
  ```
  
# SSO login

The addon service will generate credentials with the method below:
```
timestamp := time.Now().Unix()
s := sha1.New()
s.Write([]byte(fmt.Sprintf("%s:%s:%d", appSlug, addonConfig.SSOSecret, timestamp)))
token := s.Sum(nil)

c.Response().Header().Add("bitrise-sso-timestamp", fmt.Sprintf("%d", timestamp))
c.Response().Header().Add("bitrise-sso-token", fmt.Sprintf("%x", token))
c.Response().Header().Add("bitrise-sso-x-action", fmt.Sprintf("%s", fmt.Sprintf("%s/login", addonConfig.Host)))
```

and will respond with header fields(`bitrise-sso-timestamp`, `bitrise-sso-token`, `bitrise-sso-x-action`) that include those data. This is a communication between core and addon service. After the core received the data in the header, it will send a post form to the addon itself as the following:

```
method: post
action: bitrise-sso-x-action?build_slug=build_slug

fields
timestamp: bitrise-sso-timestamp
token: bitrise-sso-token
app_slug: the-appslug
```

---

> Login post-form handler.

**Method**: POST

**URL**: `/login`

**PARAMS**:

- `build_slug` (Not required to handle, it is sent by Bitrise core but it is up to the addon that it want to use for redirection or not.)

**FORMVALUES**:

- `timestamp`
- `token`
- `app_slug`
