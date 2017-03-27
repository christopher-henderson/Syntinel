from os import path, stat, remove
from shutil import copyfile, move

from django.conf import settings
from django.db import models

from syntinel.executor_bindings import ExecutorBindings


class Test(models.Model):

    name = models.CharField(max_length=2 ** 16)
    project = models.ForeignKey('Project', on_delete=models.CASCADE, null=True, related_name="tests")
    script = models.CharField(max_length=2 ** 16)
    dockerfile = models.CharField(max_length=2 ** 16)
    environmentVariables = models.CharField(max_length=2 ** 16, null=True)
    # Integer so that we can have granularity of test health.
    # Similar to Jenkin's "red, yellow, green"
    health = models.IntegerField(default=100)
    interval = models.IntegerField(null=True)


class TestRun(models.Model):

    test = models.ForeignKey('Test')
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
