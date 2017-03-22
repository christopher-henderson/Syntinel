import requests
import logging
from django.conf import settings

logger = logging.getLogger(__name__)


class ExecutorBindings(object):

    @staticmethod
    def run(id, testID, dockerfile, script, environmentVariables):
        print("Firing off to load balancer.")
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
