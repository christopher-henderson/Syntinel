from rest_framework.response import Response
from rest_framework.generics import (
    RetrieveUpdateDestroyAPIView,
    CreateAPIView,
    ListAPIView)

from rest_framework import mixins
from rest_framework import generics

from syntinel.models import (
    Test,
    Project,
    TestRun,
    Executor)
from .serializers import (
    TestSerializer,
    ProjectSerializer,
    TestRunSerializer,
    ExecutorSerializer)


class TestView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer


class TestListView(ListAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer


class ProjectView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Project.objects.all()
    serializer_class = ProjectSerializer


class ProjectListView(ListAPIView):

    queryset = Project.objects.all()
    serializer_class = ProjectSerializer


class TestRunView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = TestRun.objects.all()
    serializer_class = TestRunSerializer

    def post(self, request, pk):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        self.perform_create(serializer)
        headers = self.get_success_headers(serializer.data)
        testRun = serializer.instance
        return_code = testRun.run()
        return Response(serializer.data, status=return_code, headers=headers)


class TestRunListView(ListAPIView):

    queryset = TestRun.objects.all()
    serializer_class = TestRunSerializer


class ExecutorView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Executor.objects.all()
    serializer_class = ExecutorSerializer


class ExecutorListView(ListAPIView):

    queryset = Executor.objects.all()
    serializer_class = ExecutorSerializer
