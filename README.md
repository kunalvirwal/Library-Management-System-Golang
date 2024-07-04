# ***Steps to run:***  

- Download all the dependencies using
  > go mod tidy

- Copy all queries one by one from `db_schema.sql` onto your mysql shell.

- Create `.env` using `sample.env` file and add all the environment variables in it. 
  > cp sample.env .env  

- Bulid the program using
  > go build -o LabsLibrary ./cmd/main.go

- Run the server.
  > ./LabsLibrary

- Visit `https://localhost:4000` 
- An admin account will always be created with credentials "admin@sdslabs.com" and password "A" if no other admin account present in database.


# Possible Future Improvements
- All the model functions return user information as a struct, the code can be made faster by returning poiters to the structs
