from django.conf.urls import url

from .views import *

urlpatterns = [
    url(r'^test/(?P<pk>[0-9]+)?/?$', TestView.as_view(), name='test'),
    url(r'^test/all/?$', TestListView.as_view(), name='test_list'),

    url(r'^project/(?P<pk>[0-9]+)?/?$', ProjectView.as_view(), name='project'),
    url(r'^project/all/?$', ProjectListView.as_view(), name='project_list'),

    url(r'^testrun/(?P<pk>[0-9]+)?/?$', TestRunView.as_view(), name='testrun'),
    url(r'^testrun/all/?$', TestRunListView.as_view(), name='testrun_list'),

    url(r'^executor/(?P<pk>[0-9]+)?/?$', ExecutorView.as_view(), name='testrun'),
    url(r'^executor/all/?$', ExecutorListView.as_view(), name='testrun_list'),
    ]
