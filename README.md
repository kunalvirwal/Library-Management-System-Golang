# ***Steps to run:***

- Download all the dependencies using
  > go mod tidy

- Create `.env` using `sample.env` file and add all the environment variables in it.
  > cp sample.env .env

- In your mysql shell create an new database.
  > create database MVCdb;

- Instead Virtually host the server using apache and setup database automatically.
  > source install.sh


- Bulid the program using
  > go build -o LabsLibrary ./cmd/main.go

- Run the server.
  > ./LabsLibrary

- Visit `http://mvc.sdslabs.local` (Try https is security error on browser)
- An admin account will always be created with credentials "admin@sdslabs.com" and password "A" if no other admin account present in database.
<hr>

- To run for development purposes (Install CompileDaemon) .
  > go get github.com/githubnemo/CompileDaemon  

  > CompileDaemon -build="go build -o LabsLibrary ./cmd/main.go" -command="./LabsLibrary"


# ***Steps to run tests***

- Setup `.env` and create a database in mysql
- Make sure the program has be executed and used once (Route `/` should have been hit atleast once, and admin details has not been changed otherwiese test will fail to compare), this will setup the database.
- Run `go test -v`

# ***Steps to Virtual Host on Ubuntu linux manually***

- Install apache2 and configure:   
  > sudo apt install apache2  

  > sudo a2enmod proxy proxy_http  

  > cd /etc/apache2/sites-available  

  > sudo nano mvc.sdslabs.local.conf   

- Add:
    ```  
    <VirtualHost *:80>
      ServerName mvc.sdslabs.local
      ServerAdmin youremailid
      ProxyPreserveHost On
      ProxyPass / http://127.0.0.1:8000/
      ProxyPassReverse / http://127.0.0.1:8000/
      TransferLog /var/log/apache2/mvc_access.log
      ErrorLog /var/log/apache2/mvc_error.log
    </VirtualHost>
    ```
  > sudo a2ensite mvc.sdslabs.local.conf  

- Add `127.0.0.1	mvc.sdslabs.local` to `/etc/hosts`  

  > sudo a2dissite 000-default.conf  

  > sudo apache2ctl configtest   

  > sudo systemctl restart apache2  

  > sudo systemctl status apache2  

Check `mvc.sdslabs.local` on your browser  
<hr>  


###  Login Page
<p align=center>
<img src="ReadmeAssets/login.png">
<p/>

###  Admin Dashboard
<p align=center>
<img src="ReadmeAssets/adminDash.png">
<p/>

###  Books Catalog
<p align=center>
<img src="ReadmeAssets/bookCatalog.png">
<p/>

<p align=center>
Created by Kunal with ❤️ 
</p>