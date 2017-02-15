from os import path, stat, remove
from shutil import copyfile, move

from django.conf import settings
from django.db import models

from syntinel.executor_bindings import ExecutorBindings


class Test(models.Model):

    name = models.CharField(max_length=2 ** 16)
    suite = models.ForeignKey('Suite', on_delete=models.CASCADE, null=True)
    script = models.CharField(max_length=2 ** 16)
    dockerfile = models.CharField(max_length=2 ** 16)
    environmentVariables = models.CharField(max_length=2 ** 16, null=True)
    # Integer so that we can have granularity of test health.
    # Similar to Jenkin's "red, yellow, green"
    health = models.IntegerField(default=100)


class TestRun(models.Model):

    test = models.ForeignKey('Test')
    log = models.CharField(max_length=2 ** 16, default="")
    successful = models.NullBooleanField()

    def run(self):
        ExecutorBindings.run(
            testID=self.test.id,
            dockerfile=self.test.dockerfile,
            script=self.test.script,
            environmentVariables=self.test.environmentVariables
            )


class Suite(models.Model):

    name = models.CharField(max_length=2 ** 16)
