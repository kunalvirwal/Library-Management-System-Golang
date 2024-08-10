#!/bin/bash
path=$(pwd)
echo -e "This script will install and configure apache2, create a virtual host and help in creation of database."
echo -e "Would you like to run the installation script (Y/N)?."
read -r reply 

if [ "$reply" = "Y" ] || [ "$reply" = "y" ]
then
    sudo apt update
    sudo apt upgrade
    sudo apt install apache2
    sudo a2enmod proxy proxy_http
    cd /etc/apache2/sites-available
else
    exit 0
fi

if [ -f "mvc.sdslabs.local.conf" ]
then
    echo "Configuration File already present"
else
    echo -e "Please enter your email ID:"
    read -r emailID

    if [ -z "$emailID" ]
	then
		echo "Emailid cannot be empty"
        exit 1
	fi
    
    sudo tee -a mvc.sdslabs.local.conf > /dev/null <<EOL
<VirtualHost *:80>
	ServerName mvc.sdslabs.local
    ServerAdmin $emailID
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:8000/
    ProxyPassReverse / http://127.0.0.1:8000/
    TransferLog /var/log/apache2/mvc_access.log
    ErrorLog /var/log/apache2/mvc_error.log
</VirtualHost>
EOL
    sudo a2ensite mvc.sdslabs.local.conf
    sudo tee -a /etc/hosts > /dev/null <<EOL
127.0.0.1        mvc.sdslabs.local
EOL
    sudo a2dissite 000-default.conf
    sudo apache2ctl configtest
    sudo systemctl restart apache2
    sudo systemctl status apache2

fi

cd $path
echo -e "Would you like to setup mysq database? (Note: mysql service should be preinstalled and running) (y/n)"
read -r reply

if [ "$reply" = "Y" ] || [ "$reply" = "y" ]
then
    echo -e "Please enter your mysql username:"
    read -r sqlUser 
    echo -e "Please enter the mysql password for '$sqlUser':"
    read -r -s sqlPassword 
    mysql -u $sqlUser -p$sqlPassword -e "DROP DATABASE IF EXISTS MVCdb; CREATE DATABASE MVCdb;"
else
	echo -e "Skipping mysql database setup\n"
	exit 0
fi	
echo -e "Ready to Run \n"
