#!/usr/bin/env bash

cd /opt
source venv3/bin/activate
cd Syntinel/compozed
rm db.sqlite3
rm -rf syntinel/migrations/*
python manage.py makemigrations syntinel
python manage.py migrate syntinel
python manage.py makemigrations
python manage.py migrate
cd ../syntinel_executor
./make
/usr/bin/supervisord -c /etc/supervisord.conf
