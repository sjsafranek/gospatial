from django.contrib import admin

# Register your models here.
from .models import Layer
# from .models import Feature
from .models import Baselayer
from .models import GeoAPI
from .models import GeoAPIKey

admin.site.register(Layer)
admin.site.register(Baselayer)
admin.site.register(GeoAPI)
admin.site.register(GeoAPIKey)
