# Import the utility functions from the URL handling library
from django.conf.urls import patterns, include, url
# Import reverse_lazy method for reversing names to URLs
from django.core.urlresolvers import reverse_lazy
from django.contrib.auth import views as auth_views
#from . import views
from WallOfDenial import views

urlpatterns = [
    # Map the 'django.contrib.auth.views.login' view to the /login/ URL.
    # The additional parameters to the view are passed via the 3rd argument which is
    # a dictionary of various parameters like the name of the template to be
    # used by the view.
    url(r'^$', auth_views.login, { "template_name" : "login.html", }),
    url(r'^login/$', auth_views.login, { "template_name" : "login.html", }),
    url(r'^logout/$', auth_views.logout, { "next_page" : '/login'}),
    url(r'^help/$', views.help),
    
    # Map the 'django.contrib.auth.views.logout' view to the /logout/ URL.
    # Pass additional parameters to the view like the page to show after logout
    # via a dictionary used as the 3rd argument.
    # url(r'^home/$', views.home),
    url(r'^management/$', views.management),
    # url(r'^tracking/$', views.tracking),
    url(r'^create_layer/$', views.create_layer),
    url(r'^delete_layer/$', views.delete_layer),
    url(r'^create_baselayer/$', views.create_baselayer),
    # url(r'^delete_baselayer/$', views.delete_baselayer),
]
