# URL-Shortener with rate limiter.

We will have our user interacting  with the project and we will have a server running with golang & go server at PORT:3000.  
The go server will interact with the REDIS server which is our database at PORT:6379.

## High Level Design.
Three major components:
1. api  
    a. Database     : logic to connect to the database.  
    b. Helpers      : will have some functions to help the routes.  
    c. Routes       : routes on which we will post and shorten our URL.  
    d. Dockerfile   : to have golang specific code.  
2. data
3. db  
    a. Dockerfile   : Going to contain some container for our redis database.  


## Few commands to run
to install all the dependencies go to "api" directory and execute command: 'go mod tidy'  
The above command will intall all the dependencies.


## Command to run the project:
docker-compose up -d
