from rest_framework.response import Response
from rest_framework.generics import (
    RetrieveUpdateDestroyAPIView,
    CreateAPIView,
    ListAPIView)

from rest_framework import mixins
from rest_framework import generics

from syntinel.models import (
    Test,
    Suite,
    TestRun)
from .serializers import (
    TestSerializer,
    SuiteSerializer,
    TestRunSerializer)


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
