from django.conf.urls import url

from .views import *

urlpatterns = [
    url(r'^docker/(?P<pk>[0-9]+)?/?$', DockerView.as_view(), name='docker')
    ]
