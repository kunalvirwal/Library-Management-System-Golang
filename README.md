# ***Steps to run:***  

- Download all the dependencies using
  > go mod tidy

- Create `.env` using `sample.env` file and add all the environment variables in it. 
  > cp sample.env .env  
  
- In your mysql shell create an new database
  > create database MVCdb;

- Bulid the program using
  > go build -o LabsLibrary ./cmd/main.go

- Run the server.
  > ./LabsLibrary

- Visit `https://localhost:4000` 
- An admin account will always be created with credentials "admin@sdslabs.com" and password "A" if no other admin account present in database.


