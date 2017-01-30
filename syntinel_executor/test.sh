# curl -X POST http://localhost:8000/docker/1 --data-binary @/Users/chris/Downloads/centos.iso

# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"
# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"
# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"

curl -X POST http://192.168.1.2:8000/script/1 -d "#!/usr/bin/env bash
git clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover"
curl -X POST http://192.168.1.2:8000/docker/1 -d "
FROM docker.io/centos

MAINTAINER Christopher Henderson

RUN yum install -y go git wget
COPY script.sh \$HOME/script.sh
CMD chmod +x script.sh && ./script.sh
"
echo "Registering Test"
curl -X POST http://192.168.1.2:8000/test/1?dockerID=1\&scriptID=1
curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=4

# curl -X DELETE http://192.168.1.2:8000/test/run?testID=1\&testRunID=2

# curl -X POST http://192.168.1.2:8000/test/2?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/3?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/4?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=4\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=4\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=4\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/5?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=5\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=5\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=5\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/6?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=6\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=6\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=6\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/7?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=7\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=7\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=7\&testRunID=4
#
# curl -X POST http://192.168.1.2:8000/test/8?dockerID=1\&scriptID=1
# curl -X POST http://192.168.1.2:8000/test/run?testID=8\&testRunID=2
# curl -X POST http://192.168.1.2:8000/test/run?testID=8\&testRunID=3
# curl -X POST http://192.168.1.2:8000/test/run?testID=8\&testRunID=4
# sleep 2s
# curl -X DELETE http://localhost:8000/test/run?testID=1\&testRunID=2
# curl -X GET http://localhost:8000/test/run?testID=1\&testRunID=2
