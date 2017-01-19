# curl -X POST http://localhost:8000/docker/1 --data-binary @/Users/chris/Downloads/centos.iso

# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"
# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"
# curl -X POST http://localhost:8000/docker/1 -F "docker=@/Users/chris/Downloads/centos.iso"

curl -X POST http://192.168.1.2:8000/script/1 -d "#!/usr/bin/env bash
git clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover"
curl -X POST http://192.168.1.2:8000/docker/1 -F "docker=@/Users/chris/Documents/ASU/Syntinel/syntinel_executor/Dockerfile"
curl -X POST http://192.168.1.2:8000/test/1?dockerID=1\&scriptID=1
curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=2
curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=3
curl -X POST http://192.168.1.2:8000/test/run?testID=1\&testRunID=4

curl -X POST http://192.168.1.2:8000/test/2?dockerID=1\&scriptID=1
curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=2
curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=3
curl -X POST http://192.168.1.2:8000/test/run?testID=2\&testRunID=4

curl -X POST http://192.168.1.2:8000/test/3?dockerID=1\&scriptID=1
curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=2
curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=3
curl -X POST http://192.168.1.2:8000/test/run?testID=3\&testRunID=4
# sleep 2s
# curl -X DELETE http://localhost:8000/test/run?testID=1\&testRunID=2
# curl -X GET http://localhost:8000/test/run?testID=1\&testRunID=2
