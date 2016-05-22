from django.db import models

# Create your models here.
from django.contrib.auth.models import Group

class Layer(models.Model):
    uuid = models.CharField(max_length=50)
    name = models.CharField(max_length=100)
    server = models.CharField(max_length=150)
    apikey = models.CharField(max_length=150)
    owner = models.ForeignKey(Group)

class Baselayer(models.Model):
    url = models.CharField(max_length=150)
    name = models.CharField(max_length=150)
    owner = models.ForeignKey(Group)

class GeoAPI(models.Model):
    url = models.CharField(max_length=150)
    authkey = models.CharField(max_length=150)

class GeoAPIKey(models.Model):
    url = models.CharField(max_length=150)
    apikey = models.CharField(max_length=150)
    owner = models.ForeignKey(Group)