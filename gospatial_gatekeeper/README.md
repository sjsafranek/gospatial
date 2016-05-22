# findauth


## Quick Start

### Install requirements::
	$ sudo pip3 install django
	$ sudo pip3 install jinja2

* requires django 1.9 =<

### Setup database::
	$ python3 manage.py migrate
	$ python3 manage.py makemigrations

### Create superuser::
	$ python3 manage.py createsuperuser

### Launch server::
	$ python3 manage.py runserver

### Setup first group

	- http://localhost:8000/admin/
	- login with superuser
	- http://localhost:8000/admin/auth/group/
	- add group
	- http://localhost:8000/admin/auth/user/1/change/
	- add user to group
	- http://localhost:8000/admin/WallOfDenial/geoapi/add/
	- Create geoapi
	- assign geoapi to group
	- http://localhost:8000/
	- login to main app


