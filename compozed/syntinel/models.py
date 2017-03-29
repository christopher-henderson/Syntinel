from __future__ import division

import logging

from os import path, stat, remove
from shutil import copyfile, move

from django.conf import settings
from django.db import models

from syntinel.executor_bindings import ExecutorBindings

logger = logging.getLogger("models")


class Test(models.Model):

    NUMBER_OF_TESTS_CONSIDERED_IN_HEALTH = 10

    name = models.CharField(max_length=2 ** 16)
    project = models.ForeignKey('Project', on_delete=models.CASCADE, null=True, related_name="tests")
    script = models.CharField(max_length=2 ** 16)
    dockerfile = models.CharField(max_length=2 ** 16)
    environmentVariables = models.CharField(max_length=2 ** 16, null=True)
    # Integer so that we can have granularity of test health.
    # Similar to Jenkin's "red, yellow, green"
    health = models.IntegerField(default=100)
    interval = models.IntegerField(null=True)

    def update_health(self, successful):
        test_runs = self.test_runs.order_by('id').reverse()[:self.NUMBER_OF_TESTS_CONSIDERED_IN_HEALTH]
        total = 1 if successful else 0
        total += sum(1 for t in test_runs if t.successful)
        self.health = total / (len(test_runs)) * 100
        self.save()


class TestRun(models.Model):

    test = models.ForeignKey('Test', related_name="test_runs")
    log = models.CharField(max_length=2 ** 16, default="")
    error = models.CharField(max_length=2 ** 16, default="")
    status = models.IntegerField(null=True)
    successful = models.NullBooleanField()

    def run(self):
        ExecutorBindings.run(
            id=self.id,
            testID=self.test.id,
            dockerfile=self.test.dockerfile,
            script=self.test.script,
            environmentVariables=self.test.environmentVariables
            )

    def finalize():
        self.successful = self.status is 0


class Project(models.Model):

    name = models.CharField(max_length=2 ** 16)


class Executor(models.Model):

    hostName = models.CharField(max_length=256)
    # Max port is 65535, five characters long.
    port = models.CharField(max_length=5)
    Scheme = models.CharField(max_length=5, default="http")
