# nopwd

A simple no password library for go for doing no password style login. This library only concerns itself with logging in. Not validation of APIs after you are logged in.

```go
// What secret should tokens be validated with?
secret := "choose something random"
// How long should login links last?
ttl := 10 

noPwd := NewNoPwd("https://foo.com",secret,ttl)

func sendCodeToEmail(email string) error {
   // create a login link for an email (e.g https://foo.com/?code=ABSDIMOIAd... )
   loginLink := noPwd.GenerateCodeLink(email)
   
   // send an login with whatever tech you use for sending emails (mailgun, etc.)
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
  code := c.QueryParam("code")
  valid, err := noPwd.ValidateCode(code)
  if valid != true || err != nil {
    return c.String(http.StatusUnauthorized, "Failed to validate code")
  }
  
  // Once validated, send them whatever authorization token you want them to use for your api
  return c.String(http.StatusOK, "...")
})

...
````
