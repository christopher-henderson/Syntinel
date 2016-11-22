from django.test import TestCase

from .api.v1.models import Docker


class DockerTestCase(TestCase):

    def setUp(self):
        Docker.objects.create(name='JUnit', platform='RHEL7')
        Docker.objects.create(name='Selenium', platform='Ubuntu')

    def test_docker_get_by_platform(self):
        rhel = Docker.objects.get(platform='RHEL7')
        ubuntu = Docker.objects.get(platform='Ubuntu')
        self.assertEqual(rhel.name, 'JUnit')
        self.assertEqual(ubuntu.name, 'Selenium')

    def test_docker_get_by_name(self):
        junit = Docker.objects.get(name='JUnit')
        selenium = Docker.objects.get(name='Selenium')
        self.assertEqual(junit.platform, 'RHEL7')
        self.assertEqual(selenium.platform, 'Ubuntu')
