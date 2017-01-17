import json
import os

from django.test import TestCase
from django.test import Client

from .api.v1.models import Docker, Script

# !!!!!!!!!!!!!!
# Do NOT delete this. This tells the Script entity to save testing scripts
# to /tmp rather than the ACTUAL test script directory. If you delete this,
# then tests WILL overwrite real data.
# !!!!!!!!!!!!!!
try:
    os.mkdir('/tmp/scripts')
except Exception:
    pass
Script.scripts = '/tmp/scripts'
#######################

client = Client()


class DockerTestCase(TestCase):

    def setUp(self):
        Docker.objects.create(name='JUnit', platform='RHEL7')
        Docker.objects.create(name='Selenium', platform='Ubuntu')
        Docker.objects.create(name='Integration', platform='CoreOS')
        Docker.objects.create(name='Smoke', platform='CoreOS')

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

    def test_docker_filter(self):
        objs = Docker.objects.filter(platform='CoreOS')
        self.assertGreater(len(objs), 1)

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

    def test_docker_client_all(self):
        response = client.get('/api/v1/docker/all')
        self.assertEqual(response.status_code, 200)
        dockers = json.loads(response.content.decode())
        self.assertEqual(len(Docker.objects.all()), len(dockers))
        for docker in dockers:
            db_obj = Docker.objects.get(id=docker['id'])
            self.assertEqual(db_obj.id, docker['id'])
            self.assertEqual(db_obj.name, docker['name'])
            self.assertEqual(db_obj.platform, docker['platform'])


class ScriptTestCase(TestCase):

    def setUp(self):
        s1 = Script.objects.create()
        s1.content = '#/usr/bin/env python'
        s2 = Script.objects.create()
        s2.content = '#/usr/bin/env bash'

    def test_script_client_post(self):
        script = 'echo "Hello World."'
        response = client.post('/api/v1/script/', data={'content': script})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())
        self.assertEqual(obj['content'], script)
        db_obj = Script.objects.get(id=obj['id'])
        self.assertEqual(obj['content'], db_obj.content)

    def test_script_client_get(self):
        script = 'echo "Hello World."'
        response = client.post('/api/v1/script/', data={'content': script})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())

        response = client.get('/api/v1/script/{ID}'.format(ID=obj['id']))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(obj['content'], script)
        db_obj = Script.objects.get(id=obj['id'])
        self.assertEqual(obj['content'], db_obj.content)

    def test_script_client_patch(self):
        script = 'echo "Hello World."'
        response = client.post('/api/v1/script/', data={'content': script})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())

        new_script = 'javac lotsOfWords.java'
        response = client.patch(
            '/api/v1/script/{ID}'.format(ID=obj['id']),
            data=json.dumps({'content': new_script}),
            content_type='application/json')
        self.assertEqual(response.status_code, 200)
        new_obj = json.loads(response.content.decode())
        self.assertEqual(new_obj['content'], new_script)
        db_obj = Script.objects.get(id=new_obj['id'])
        self.assertEqual(new_obj['content'], db_obj.content)

    def test_script_client_put(self):
        script = 'echo "Hello World."'
        response = client.post('/api/v1/script/', data={'content': script})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())

        new_script = 'javac lotsOfWords.java'
        response = client.put(
            '/api/v1/script/{ID}'.format(ID=obj['id']),
            data=json.dumps({'content': new_script}),
            content_type='application/json')
        self.assertEqual(response.status_code, 200)
        new_obj = json.loads(response.content.decode())
        self.assertEqual(new_obj['content'], new_script)
        db_obj = Script.objects.get(id=new_obj['id'])
        self.assertEqual(new_obj['content'], db_obj.content)

    def test_script_client_delete(self):
        script = 'echo "Hello World."'
        response = client.post('/api/v1/script/', data={'content': script})
        self.assertEqual(response.status_code, 201)
        obj = json.loads(response.content.decode())

        response = client.delete('/api/v1/script/{ID}'.format(ID=obj['id']))
        self.assertEqual(response.status_code, 204)
        try:
            Script.objects.get(id=obj['id'])
        except Exception:
            # This is what we want.
            pass
        else:
            raise AssertionError("Failed to delete script entity from the DB.")
        try:
            with open(os.path.join(Script.scripts, str(obj['id'])), 'r') as s:
                pass
        except:
            pass
        else:
            raise AssertionError('Failed to delete script from the FS.')

    def test_script_client_all(self):
        response = client.get('/api/v1/script/all')
        self.assertEqual(response.status_code, 200)
        objs = json.loads(response.content.decode())
        self.assertGreater(len(objs), 1)
