from os import path

from django.conf import settings
from django.db import models


class Docker(models.Model):

    name = models.CharField(max_length=2 ** 16)
    platform = models.CharField(max_length=2 ** 16)


class Test(models.Model):

    name = models.CharField(max_length=2 ** 16)
    suite = models.ForeignKey('Suite', on_delete=models.CASCADE)
    scripts = models.ManyToManyField('Script')
    # Integer so that we can have granularity of test health.
    # Similar to Jenkin's "red, yellow, green"
    health = models.IntegerField()


class TestRun(models.Model):

    test = models.ForeignKey('Test')
    log = models.CharField(max_length=2 ** 16)
    successful = models.BooleanField()


class Suite(models.Model):

    name = models.CharField(max_length=2 ** 16)
    docker = models.ForeignKey('Docker')


class Script(models.Model):

    # Path to script stored on the filesystem. This is to avoid varchar
    # limits in our database (e.g. 65k in mysql).
    path = models.CharField(max_length=2 ** 16)

    # Absolute path to the static scripts directory.
    scripts = path.join(settings.BASE_DIR, 'syntinel', 'static', 'scripts')

    def __init__(self, *args, **kwargs):
        self._unsaved_content = None
        if args and not kwargs:
            super(Script, self).__init__(*args, **kwargs)
        else:
            self._unsaved_content = kwargs.get('content')
            super(Script, self).__init__(*args, {})

    def save(self):
        pass
