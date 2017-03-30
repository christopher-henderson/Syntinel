import logging
import json
import datetime

from django.core import cache
from django.core import exceptions
from django.utils import timezone

from django_filters.rest_framework import DjangoFilterBackend

from rest_framework.response import Response
from rest_framework.generics import (
    RetrieveUpdateDestroyAPIView,
    CreateAPIView,
    ListAPIView,
    UpdateAPIView)
from rest_framework.status import (
    HTTP_200_OK,
    HTTP_400_BAD_REQUEST,
    HTTP_201_CREATED,
    )

from rest_framework import mixins
from rest_framework import generics

from syntinel.consumers import LogCache

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

from syntinel.executor_bindings import ExecutorBindings

logger = logging.getLogger("views")


class TestView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer

    def post(self, request, pk):
        try:
            wasScheduled = Test.objects.get(id=pk).interval is not None
        except Test.DoesNotExist:
            wasScheduled = False
        serializer = TestSerializer(data=request.data)
        if serializer.is_valid():
            serializer.save()
        else:
            return Response(serializer.errors, status=HTTP_400_BAD_REQUEST)
        test = serializer.instance
        isScheduled = test.interval is not None
        logger.debug("wasScheduled: {W}".format(W=wasScheduled))
        logger.debug("isScheduled: {I}".format(I=isScheduled))
        logger.debug(request.data)
        if not wasScheduled and isScheduled:
            loadbalancer_status = ExecutorBindings.schedule(test)
            logger.debug("Received status code {C} from load balancer while scheduling test ID {ID}".format(
                C=loadbalancer_status,
                ID=test.id
                ))
        elif wasScheduled and not isScheduled:
            loadbalancer_status = ExecutorBindings.cancel(test)
            logger.debug("Received status code {C} from load balancer while cancelling scheduling of test ID {ID}".format(
                C=loadbalancer_status,
                ID=test.id
                ))
        return Response(serializer.data, status=HTTP_201_CREATED)

    def patch(self, request, pk):
        test = Test.objects.get(id=pk)
        wasScheduled = test.interval is not None
        response = super(TestView, self).patch(request, pk)
        interval = request.data.get("interval")
        isScheduled = interval is not None
        logger.debug("in patch")
        if not wasScheduled and isScheduled:
            test.interval = interval
            status = ExecutorBindings.schedule(test)
            logger.debug("Received status code {C} from load balancer while scheduling test ID {ID}".format(
                C=status,
                ID=pk
                ))
        elif wasScheduled and not isScheduled:
            status = ExecutorBindings.cancel(test)
            logger.debug("Received status code {C} from load balancer while cancelling scheduling of test ID {ID}".format(
                C=status,
                ID=pk
                ))
        return response


class TestListView(ListAPIView):

    queryset = Test.objects.all()
    serializer_class = TestSerializer
    filter_fields = ("project",)


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

    def patch(self, request, pk):
        logger.debug("Finalizing test " + str(pk))
        logCache = LogCache.getLogCache(pk)
        logCache.finalize()
        test_run = TestRun.objects.get(id=pk)
        test_run.completed = datetime.datetime.now(test_run.started.tzinfo)
        test_run.duration = test_run.completed - test_run.started
        test_run.test.update_health(request.data.get("successful"))
        test_run.save()
        return super(TestRunView, self).patch(request, pk)


class TestRunListView(ListAPIView):

    queryset = TestRun.objects.all()
    serializer_class = TestRunSerializer
    filter_fields = ("test", "successful")


class ExecutorView(CreateAPIView, RetrieveUpdateDestroyAPIView):

    queryset = Executor.objects.all()
    serializer_class = ExecutorSerializer


class ExecutorListView(ListAPIView):

    queryset = Executor.objects.all()
    serializer_class = ExecutorSerializer
