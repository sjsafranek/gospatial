
from .models import *
from django.contrib.auth.models import User
from django.contrib.auth.models import Group

def get_users_by_group(groupname):
    group = Group.objects.get(name=groupname)
    users = group.user_set.all()
    return users

def get_baselayers_by_group(groupname):
    group = Group.objects.get(name=groupname)
    baselayers = Baselayer.objects.all().filter(owner=group)
    return baselayers

def get_layers_by_group(groupname):
    group = Group.objects.get(name=groupname)
    layers = Layer.objects.all().filter(owner=group)
    return layers

def getGeoAPIKey(groupname):
    group = Group.objects.get(name=groupname)
    engine = group.geoapikey_set.all()
    return {
        "apikey": str(engine[0].apikey),
        "address": str(engine[0].url)
    }
