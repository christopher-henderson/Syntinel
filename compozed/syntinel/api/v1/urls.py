from django.conf.urls import url

from .views import *

urlpatterns = [
    url(r'^docker/(?P<pk>[0-9]+)?/?$', DockerView.as_view(), name='docker'),
    url(r'^docker/all/?$', DockerListView.as_view(), name='docker_list'),

    url(r'^test/(?P<pk>[0-9]+)?/?$', TestView.as_view(), name='test'),
    url(r'^test/all/?$', TestListView.as_view(), name='test_list'),

    url(r'^suite/(?P<pk>[0-9]+)?/?$', SuiteView.as_view(), name='suite'),
    url(r'^suite/all/?$', SuiteListView.as_view(), name='suite_list'),

    url(r'^script/(?P<pk>[0-9]+)?/?$', ScriptView.as_view(), name='script'),
    url(r'^script/all/?$', ScriptListView.as_view(), name='script_list'),

    url(r'^testrun/(?P<pk>[0-9]+)?/?$', TestRunView.as_view(), name='testrun'),
    url(r'^testrun/all/?$', TestRunListView.as_view(), name='testrun_list'),
    ]
