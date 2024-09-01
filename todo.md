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

board onboarding flow
[] boards include information about a client
[] each client has a public key that the board sends to the server to validate against a private key
[] board per hour pings the server and asks to be assigned a machine untill its assighe\