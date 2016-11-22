from django.conf.urls import url

from .views import *

urlpatterns = [
    url(r'docker/?', DockerView.as_view())
    ]
