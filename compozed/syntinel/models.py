from django.db import models


class Docker(models.Model):

    name = models.CharField()
    image_type = models.CharField()


class Test(models.Model):

    name = models.CharField()
    suite = models.ForeignKey('Suite', on_delete=models.Cascade)
    scripts = models.ManyToManyField('Scripts')
    # Integer so that we can have granularity of test health.
    # Similar to Jenkin's "red, yellow, green"
    health = models.IntegerField()


class Suite(models.Model):

    name = models.CharField()
    docker = models.ForeignKey('Docker')


class Script(models.Model):

    # Path to script stored on the filesystem. This is to avoid varchar
    # limits in our database (e.g. 65k in mysql).
    content = models.CharField()
