from os import path, stat, remove
from shutil import copyfile, move

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
    # path = models.CharField(max_length=2 ** 16)

    # Absolute path to the static scripts directory.
    scripts = path.join(settings.BASE_DIR, 'syntinel', 'static', 'scripts')

    def delete(self):
        try:
            remove(self.script_path)
        except Exception:
            pass
        super(Script, self).delete()

    @property
    def content(self):
        c = ''
        try:
            with open(self.script_path, 'r') as script:
                c = ''.join(script.readlines())
        except FileNotFoundError:
            pass
        return c

    @content.setter
    def content(self, content):
        self.backup()
        try:
            with open(self.script_path, 'w+') as script:
                script.writelines(content)
        except Exception as error:
            settings.logger.error(error)
            self.restore()
        else:
            self.cleanup_backup()

    @property
    def script_path(self):
        return path.join(self.scripts, str(self.id))

    def backup(self):
        try:
            stat(self.script_path)
        except FileNotFoundError:
            pass
        else:
            copyfile(self.script_path, path.join('/tmp', str(self.id)))

    def restore(self):
        tmp = path.join('/tmp', str(self.id))
        try:
            stat(tmp)
        except Exception:
            pass
        else:
            move(tmp, self.script_path)

    def cleanup_backup(self):
        tmp = path.join('/tmp', str(self.id))
        try:
            stat(tmp)
        except Exception:
            pass
        else:
            remove(tmp)
