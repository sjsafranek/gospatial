# Import redirect and render shortcuts
from django.shortcuts import render
from django.shortcuts import redirect
from django.shortcuts import render_to_response

# Import reverse_lazy method for reversing names to URLs
from django.core.urlresolvers import reverse_lazy

# Import the login_required decorator which can be applied to views 
# to enforce that the user should be logged in to access the view
from django.contrib.auth.decorators import login_required

import os
import json
import requests
from .models import *
from .Conf import *
import WallOfDenial.utils as utils
from django.http import HttpResponse
from django.http import JsonResponse
from django.contrib.auth.models import User
from django.contrib.auth.models import Group
from django.contrib.auth.hashers import check_password

# My hacky way of doing Jinja2 rendering because
# Django Jinja2 GIVES NO ERROR OUTPUT WHEN SOMETHING GOES WRONG - it just quits
# so fuck that.
from jinja2 import Environment, FileSystemLoader
env = Environment(loader=FileSystemLoader('WallOfDenial/templates'))

# import the logging library
import logging
# Get an instance of a logger
logger = logging.getLogger(__name__)
'''
Logging -
    logger.debug()
    logger.info()
    logger.warning()
    logger.error()
    logger.critical()
'''

@login_required(login_url='/')
def help(request):
    if request.method == "GET":
        if request.user.is_authenticated():
            return render(request, "help.html")
        else:
            return redirect('/login')  

@login_required(login_url='/')
def management(request):
    if request.method == "GET":
        try:
            user = User.objects.get(username=request.user.username)
            group = user.groups.all()[0]
            results = {
                'username': request.user.username,
                'group': group.name,
                'users': [],
                'layers': {},
                'baselayers': {}
            }
            for user in utils.get_users_by_group(group.name):
                results['users'].append(user.username)
            for baselayer in utils.get_baselayers_by_group(group.name):
                results['baselayers'][baselayer.name] = baselayer.url
            for layer in utils.get_layers_by_group(group.name):
                results['layers'][layer.name] = { 
                    "datasource_id": layer.uuid,
                    "apiserver": layer.server,
                    "apikey": layer.apikey
                }
            return render(request, "group_management.html",results)
        except Exception as e:
            logger.error(e)
            return redirect('/error')

@login_required(login_url='/')
def create_layer(request):
    if request.method == "POST":
        user = User.objects.get(
                    username=request.user.username)
        group = user.groups.all()[0]
        print(group.name)
        apiserver = utils.getGeoAPIKey(group.name)
        params = {"apikey": apiserver['apikey']}
        req = requests.post(
                    apiserver['address'] + "/api/v1/layer", 
                    params=params)
        if req.status_code != 200:
            raise ValueError(req.text)
        res = json.loads(req.json())
        ds = res["datasource"]
        layer = Layer.objects.create(
                    name=request.POST["name"],
                    server=apiserver['address'],
                    apikey=apiserver['apikey'],
                    uuid=ds,
                    owner=group)
        layer.save()
        return redirect('/management')

@login_required(login_url='/')
def delete_layer(request):
    if request.method == "POST":
        user = User.objects.get(
                    username=request.user.username)
        group = user.groups.all()[0]
        ds = request.POST['layer']
        apiserver = utils.getGeoAPIKey(group.name)
        params = {"apikey": apiserver['apikey']}
        req = requests.delete(
                    apiserver['address'] + "/api/v1/layer/" + ds, 
                    params=params)
        if req.status_code != 200:
            raise ValueError(req.text)
        Layer.objects.all().filter(uuid=ds).delete()
        return JsonResponse(json.loads(req.json()))

@login_required(login_url='/')
def create_baselayer(request):
    if request.method == "POST":
        user = User.objects.get(
                    username=request.user.username)
        group = user.groups.all()[0]
        baselayer = Baselayer.objects.create(
                        name=request.POST["name"],
                        url=request.POST["url"],
                        owner=group)
        baselayer.save()
        return redirect('/management')


@login_required(login_url='/')
def map(request):
    if request.method == "GET":
        try:
            user = User.objects.get(username=request.user.username)
            group = user.groups.all()[0]
            results = {
                'username': request.user.username,
                'group': group.name,
                'users': [],
                'layers': {},
                'baselayers': {},
                'servers': json.dumps({
                    'gis': utils.getGeoAPIKey(group)
                })
            }
            for user in utils.get_users_by_group(group.name):
                results['users'].append(user.username)
            for baselayer in utils.get_baselayers_by_group(group.name):
                results['baselayers'][baselayer.name] = baselayer.url
            for layer in utils.get_layers_by_group(group.name):
                results['layers'][layer.name] = layer.uuid
            return render(request, "map.html",results)
        except Exception as e:
            logger.error(e)
            return redirect('/error')

