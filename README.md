# URL-Shortener with rate limiter.

We will have our user interacting  with the project and we will have a server running with golang & go server at PORT:3000. The go server will interact with the REDIS server which is our database at PORT:6379.

## High Level Design.
Three major components:
1. api
    i) Database     : logic to connect to the database.
   ii) Helpers      : will have some functions to help the routes.
  iii) Routes       : routes on which we will post and shorten our URL.
   iv) Dockerfile   : to have golang specific code.
2. data
3. db
    i) Dockerfile   : Going to contain some container for our redis database.

