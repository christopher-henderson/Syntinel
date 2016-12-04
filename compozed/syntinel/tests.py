import json

from django.test import TestCase
from django.test import Client

from .api.v1.models import Docker
from .api.v1.serializers import DockerSerializer

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

    def test_docker_get_all(self):
        objs = Docker.objects.all()
        self.assertGreater(len(objs), 0)

    def test_docker_client_post(self):
        response = client.post(
            '/api/v1/docker/', data={'name': 'Test', 'platform': 'RHEL7'})
        self.assertEqual(response.status_code, 201)

    def test_docker_client_patch(self):
        response = client.post(
            '/api/v1/docker/', data={'name': 'Test', 'platform': 'RHEL7'})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())
        response = client.patch(
            '/api/v1/docker/{ID}'.format(ID=obj['id']),
            data=json.dumps({"name": "Updated"}),
            content_type='application/json')
        self.assertEqual(response.status_code, 200)
        updated_obj = json.loads(response.content.decode())
        self.assertEqual(updated_obj['name'], 'Updated')
        db_obj = Docker.objects.get(id=updated_obj['id'])
        self.assertEqual(updated_obj['name'], db_obj.name)

    def test_docker_client_put(self):
        response = client.post(
            '/api/v1/docker/', data={'name': 'Test', 'platform': 'RHEL7'})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())
        response = client.put(
            '/api/v1/docker/{ID}'.format(ID=obj['id']),
            data=json.dumps({'name': 'Updated', 'platform': 'Totally new'}),
            content_type='application/json')
        self.assertEqual(response.status_code, 200)
        updated_obj = json.loads(response.content.decode())
        self.assertEqual(updated_obj['name'], 'Updated')
        self.assertEqual(updated_obj['platform'], 'Totally new')
        db_obj = Docker.objects.get(id=updated_obj['id'])
        self.assertEqual(updated_obj['name'], db_obj.name)
        self.assertEqual(updated_obj['platform'], db_obj.platform)

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
