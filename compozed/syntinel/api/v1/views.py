from rest_framework.response import Response
from rest_framework.generics import (
    RetrieveUpdateDestroyAPIView,
    CreateAPIView,
    ListAPIView)

from rest_framework import mixins
from rest_framework import generics

from .models import (
    Docker,
    Test,
    Suite,
    Script,
    TestRun)
from .serializers import (
    DockerSerializer,
    TestSerializer,
    SuiteSerializer,
    ScriptSerializer,
    TestRunSerializer)


class DockerView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Docker.objects.all()
    serializer_class = DockerSerializer


class DockerListView(ListAPIView):

    queryset = Docker.objects.all()
    serializer_class = DockerSerializer


class TestView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer


class TestListView(ListAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer


class SuiteView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Suite.objects.all()
    serializer_class = SuiteSerializer


class SuiteListView(ListAPIView):

    queryset = Suite.objects.all()
    serializer_class = SuiteSerializer


class ScriptView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Script.objects.all()
    serializer_class = ScriptSerializer


class ScriptListView(ListAPIView):

    queryset = Script.objects.all()
    serializer_class = ScriptSerializer


class TestRunView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = TestRun.objects.all()
    serializer_class = TestRunSerializer


class TestRunListView(ListAPIView):

    queryset = TestRun.objects.all()
    serializer_class = TestRunSerializer
