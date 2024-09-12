* [x] auth
  * [x] gen token
  * [x] validate token
  * [x] middleware
  * [x] route for user and device
  * [x] device business logic
  * [] user business logic
* [x] firmware routes
  * [x] get latest
  * [x] get file
* [x] status route
  * [x] update status
  * [x] get status
* [x] board management
  * [x] gen token
  * [x] invalidate board 
* []x client management endpoints
  * [x] add new client
  * [x] add new buildings
  * [x] add new machines
* [] metrics
  * [] requests per second
  * [] rate of change?
  * [] often times used
* [x] logging
  * [x] logging middleware 
* [] load testing
  * [] spike test
  * [] load test
* [] request validation
* [] client info route
  * []get all buildings and the machines in the building
* [x] onboarding 
  * [x] generate new client key
  * [x] encrypt client key
  * [x] compare client key to encrypted key
* [] clean up code
* [] make a general message struct
* [] add a custom error route
* [] simplify local dev 
* [] update docker file
* [] document api
* [] metrics
  * [] return how busy each hour is at each day of the week
  * [] return how busy an individual machine is based on machine id
* [] Increase reliability with context