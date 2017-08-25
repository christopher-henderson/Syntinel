from json import dumps
from os.path import dirname

pwd = dirname(__file__)

with open(pwd + "/Dockerfile", "r") as f:
    dockerfile = ''.join(line for line in f.readlines())

with open(pwd + "/script", "r") as f:
    script = ''.join(line for line in f.readlines())

print(dumps(
    {
        "id": 1,
        "testID": 1,
        "dockerfile": dockerfile,
        "script": script,
        "environmentVariables": "a=b"
        }
    ))
