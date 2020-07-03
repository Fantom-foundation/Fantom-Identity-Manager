# Login App & Identity Manager

An application implementing the login and register workflows and storing user data.
Application will be accessed by OpenID server to provide user authentication.

Application is build on top of [Authboss](https://github.com/volatiletech/authboss) and user flows are based on [authboss sample](https://github.com/volatiletech/authboss-sample).

## Configuration

Required parameters are shown in `run` operation in `Makefile`

Application is configured to use in-memory store of user data.
Users are loaded on each application run from `users.example.json`.

OpenId server should be run with this application for proper function.