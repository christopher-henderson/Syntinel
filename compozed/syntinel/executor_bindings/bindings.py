import requests
import logging
from django.conf import settings

logger = logging.getLogger("bindings")


class ExecutorBindings(object):

    @staticmethod
    def run(id, testID, dockerfile, script, environmentVariables):
        data = {
            "ID": id,
            "testID": testID,
            "dockerfile": dockerfile,
            "script": script,
            "environmentVariables": environmentVariables
            }
        url = "http://{URL}/test/run".format(URL=settings.LOAD_BALANCER)
        response = requests.post(
            url=url,
            json=data
            )
        return response.status_code

    @staticmethod
    def schedule(test):
        data = {
            "TestID": test.id,
            "Interval": test.interval
            }
        url = "http://{URL}/schedule".format(URL=settings.LOAD_BALANCER)
        response = requests.post(
            url=url,
            json=data
            )
        logger.debug(response.text)
        return response.status_code

    @staticmethod
    def cancel(test):
        data = {"TestID": test.id}
        url = "http://{URL}/cancel".format(URL=settings.LOAD_BALANCER)
        response = requests.post(
            url=url,
            json=data
            )
        logger.debug(response.text)
        return response.status_code
