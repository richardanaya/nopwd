# nopwd

A simple no password library for go for generating login links to send to your user with a login code that can obtain an api code. Codes are:

* validated with HMAC256 secret
* validates the issuer website is the same
* validates the code is not validated beyond the time to live (TTL) you specify (implying the link could be used multiple times before the TTL expires)

This library can also generate API codee, login codes cannot be used as api codes as well as api codes cannot be used as login.

## API

* `GenerateLoginLink(email,ttl) (string,error)`
* `GenerateLoginCode(email,ttl) (string,error)`
* `ValidateLoginCode(code) (bool,string,error)`
* `GenerateAPICode(email,ttl) (string,error)`
* `ValidateAPICode(code) (bool,string,error)`

**Note: this library is not enough for production grade login system, but it might be good for experiments. Consider more issues such as rate limitting, blacklisting, asking for a new code before the old one expires. If you are in doubt, look at the source of this project. It's very minimal.**

```go
// What secret should tokens be validated with?
var secret = "choose something random"
// Global no password
var noPwd = NewNoPwd("https://foo.com",secret)

func sendCodeToEmail(email string) error {
   // create a login link for an email (e.g https://foo.com/?login_code=ABSDIMOIAd... )
   // that lasts for 10 minutes
   loginLink := noPwd.GenerateLoginLink(email, 10)
   
   // send login link with whatever tech you use for sending emails (mailgun, etc.)
   ...
}

...

e := echo.New()
e.POST("/send_code", func(c echo.Context) error {
  email := c.QueryParam("email")
  sendCodeToEmail(email)
  return c.String(http.StatusOK, "OK")
})

e.POST("/login", func(c echo.Context) error {
  code := c.QueryParam("login_code")
  valid, email, err := noPwd.ValidateLoginCode(code)
  if valid != true || err != nil {
    return c.String(http.StatusUnauthorized, "Failed to validate code")
  }
  
  // Once validated, send them whatever authorization token you 
  // want them to use for your api OR you can use api codes generated by NoPwd

  // Create an API code for an email that lasts 
  // a month (defined in minutes)
  api_code = noPwd.GenerateAPICode(email,43800) 

  return c.JSONBlob(http.StatusOK, []byte(`{
    api_code: `+api_code+`
  }`))
})


e.POST("/some_api", func(c echo.Context) error {
  code := c.Request().Headers.Get("MyAuthHeader")
  valid, email, err := noPwd.ValidateAPICode(code)
  if valid != true || err != nil {
    return c.String(http.StatusUnauthorized, "Failed to validate code")
  }
  
  // do something
})
...
````
