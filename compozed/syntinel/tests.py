from django.test import TestCase
from django.test import Client

from .api.v1.models import Docker

client = Client()


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

    def test_docker_client_post(self):
        response = client.post('/api/v1/docker/', data={'name':'Test', 'platform':'RHEL7'})
        self.assertEqual(response.status_code, 201)

    def test_docker_client_get(self):
        rhel = Docker.objects.get(platform='RHEL7')
        response = client.get('/api/v1/docker/{ID}'.format(ID=rhel.id))
        self.assertEqual(response.status_code, 200)

    def test_docker_client_delete(self):
        mock = Docker.objects.create(name='Selenium', platform='Ubuntu')
        response = client.delete('/api/v1/docker/{ID}'.format(ID=mock.id))
        self.assertEqual(response.status_code, 204)
        try:
            Docker.objects.get(pk=mock.id)
        except:
            pass
        else:
            raise AssertionError('Failed to delete Docker object.')
