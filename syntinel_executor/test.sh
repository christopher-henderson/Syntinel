curl -X POST http://localhost:9090/test/run -d '{"id": 1, "dockerfile": "FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh\n", "environmentVariables": "a=b", "script": "#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover\n", "testID": 1}'
curl -X GET http://localhost:9091/test/run?testID=1\&testRunID=1
curl -X POST http://localhost:9093/test/run -d '{"id": 69, "dockerfile": "FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh\n", "environmentVariables": "a=b", "script": "#!/usr/bin/env bash\ngit clone https://github.com/qbecker/TestCancelTest.git && cd TestCancelTest/GoBeInGoodHands && go test . -v -cover\n", "testID": 1}'
curl -X DELETE http://localhost:9091/test/run?testID=1\&testRunID=2





curl -X POST http://localhost:9093/register -d'[{"hostName": "localhost", "port": "9090", "Scheme": "http"}]'
