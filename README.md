# Bitrise sample add-on service

This project imitates the Bitrise add-on interface, you have to implement, when you're developing an add-on for Bitrise. You can run it locally and test your service in action.

The only requirement to develop a 3rd party Bitrise add-on is that the server must implement the Bitrise add-on interface. This contains 4 endpoint you have to handle a `/provision` endpoint with the methods: `POST`, `PUT`, `DELETE` and a `/login` endpoint with `POST` method. The `/provision` endpoints create, update and destroy the connection between your add-on and Bitrise. At the provisioning process you can specify the variables, which you want Bitrise to send to VMs for a build (e.g. access token, so the VM can authenticate to the add-on or the URL of your add-on). You can have as many addon server logic as you need, however the provision endpoints must use different authentication than the others, it must be the shared token that you will give to Bitrise, so your add-on can be added to the add-on service.

We will need these infos when you want to add a new addon: https://github.com/bitrise-io/bitrise-addon-service/blob/master/_template/sample-addon.yml

# Authentication

All 3 provisioning endpoints require authentication. This must be the shared token from the shared configuration. The header key must be `Authentication` and the value can be the shared token itself. Beside of the 3 provision handler, all of your other endpoints must use different authentication.

Example authentication header key-value pair:
 - `"Authentication" : "bitrise-shared-token"`

---

# Provision of a new app

The server creates a new record or updates an existing one with the app slug, to store provision state of the app. Also store a unique token for the app slug that will be used for the requests that are from a Bitrise build and calls your add-on server. Also store the received plan, so you can have a service that can use specified parameters/limits by the plan. Finally send back the list of environment variables that will be exported in all of the builds on Bitrise for the app.

**Method**: POST

**URL**: `/provision`

**Content**:

```
{
    "plan": "free",
    "app_slug": "app-slug-123",
    "api_token": "public-API-token",
    "app_title": "My awesome add-on"
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

Overwrite the plan that you saved already with the one that is in this request. This way Bitrise can update your addon if there was a plan change for any reason.

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

Delete the app's provisioned state, so the calls are pointed to this service will be rejected in the Bitrise build.

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

The addon service will generate credentials with an implementation similar to the code snippet below:

```
timestamp := time.Now().Unix()
s := sha256.New()
s.Write([]byte(fmt.Sprintf("%s:%s:%d", appSlug, addonConfig.SSOSecret, timestamp)))
token := s.Sum(nil)
tokenStr := fmt.Sprintf("sha256-%x", token)

c.Response().Header().Add("bitrise-sso-timestamp", fmt.Sprintf("%d", timestamp))
c.Response().Header().Add("bitrise-sso-token", fmt.Sprintf("%x", tokenStr))
c.Response().Header().Add("bitrise-sso-x-action", fmt.Sprintf("%s", fmt.Sprintf("%s/login", addonConfig.Host)))
```

and will respond with header fields(`bitrise-sso-timestamp`, `bitrise-sso-token`, `bitrise-sso-x-action`) which include those data. This is a communication between Bitrise core and add-on service. After the core received the data in the header, it will send a post form to the add-on itself as the following:

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

- `build_slug` (Not required to handle, it is sent by Bitrise core but it is up to the add-on that it want to use for redirection or not.)
- `app_title` provided by Bitrise, so Beam can show the Bitrise app's name, where the add-on is provisioned

**FORMVALUES**:

- `timestamp`
- `token`
- `app_slug`
