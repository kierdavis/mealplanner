#!/bin/sh

cdir=`pwd`
cd /tmp

bits=32
if [ _`uname -m` = _x86_64 ]; then
    bits=64
fi

wget -O mealplanner-server "https://github.com/kierdavis/mealplanner/raw/releases/v1.0.1/mealplanner_1.0.1_linux_$bits"
wget -O resources.zip "https://raw.github.com/kierdavis/mealplanner/releases/v1.0.1/resources.zip"

echo "Extracting resources."
unzip resources.zip
rm resources.zip

echo "Installing mealplanner-server to /usr/local/bin."
echo "You may need to enter your password to continue."
sudo mkdir -p /usr/local/bin
sudo mv mealplanner-server /usr/local/bin
sudo chmod 0755 /usr/local/bin/mealplanner-server

echo "Installing resources to /usr/local/share/mealplanner."
sudo mkdir -p /usr/local/share/mealplanner
sudo mv resources /usr/local/share/mealplanner

echo
echo "You will now be asked to enter the settings for your database."
echo "If you press Enter without typing anything the contents of the square brackets will be used as the default."

echo -n "Hostname: [localhost] "
read dbhostname
echo -n "Port: [3306] "
read dbport
echo -n "Username: [mealplanner] "
read dbusername
echo -n "Password: [] "
read dbpassword
echo -n "Database name: [mealplanner] "
read dbname

if [ -z "$dbhostname" ]; then
    dbhostname="localhost"
fi

if [ -z "$dbport" ]; then
    dbport="3306"
fi

if [ -z "$dbusername" ]; then
    dbusername="mealplanner"
fi

if [ -z "$dbname" ]; then
    dbname="mealplanner"
fi

dbsource="$dbusername:$dbpassword@tcp($dbhostname:$dbport)/$dbname"

export MPDBSOURCE="$dbsource"
export MPRESDIR="/usr/local/share/mealplanner/resources"

echo
echo "The installation is now complete."
echo "Please place the following lines into a startup shell script, such as your .bashrc:"
echo
echo "  export MPDBSOURCE='$dbsource'"
echo "  export MPRESDIR='/usr/local/share/mealplanner/resources'"
echo

cd "$cdir"
