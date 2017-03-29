---
title: Syntinel API Reference

language_tabs:
  - shell
  - javascript

toc_footers:
  - <a href='#'>Sign Up for a Developer Key</a>
  - <a href='https://github.com/tripit/slate'>Documentation Powered by Slate</a>

includes:
  - errors

search: true
---

# Introduction

Syntinel API Documentation

# Tests

## Create a Test

Creating a test requires JSON body formatted as follows:

name: String naming the test.<br>
dockerfile: String field of the full dockerfile.<br>
script: String field of the full test script.<br>
environmentVariables: Key/value pairs. E.G ```a=b, c=d```<br>

```shell
curl -X POST -H "Content-type:application/json" \
http://192.168.1.2/api/v1/test/ -d \
'{"environmentVariables": "a=b", "dockerfile": "FROM docker.io/centos\\n\\nMAINTAINER Christopher Henderson\\n\\nRUN yum install -y go git wget\\nCOPY script.sh $HOME/script.sh\\nCMD chmod +x script.sh && ./script.sh", "name": "The greatest song in the world", "script": "#!/usr/bin/env bash\\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover", "interval": 300}'
```

```javascript
var test = {
  "name":"The greatest song in the world",
  "script":"#!/usr/bin/env bash\\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover",
  "dockerfile":"FROM docker.io/centos\\n\\nMAINTAINER Christopher Henderson\\n\\nRUN yum install -y go git wget\\nCOPY script.sh $HOME/script.sh\\nCMD chmod +x script.sh && ./script.sh",
  "environmentVariables":"a=b",
  "interval": 300
}
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("POST", "http://192.168.1.2/api/v1/test/");
xmlHttp.setRequestHeader("Content-Type", "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var test = JSON.parse(xmlHttp.responseText);
    console.log(test);
  }
};
xmlHttp.send(JSON.stringify(test));
```

> The above command returns JSON structured like this:

```json
{
  "id":3,
  "name":"The greatest song in the world",
  "script":"#!/usr/bin/env bash\\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover",
  "dockerfile":"FROM docker.io/centos\\n\\nMAINTAINER Christopher Henderson\\n\\nRUN yum install -y go git wget\\nCOPY script.sh $HOME/script.sh\\nCMD chmod +x script.sh && ./script.sh",
  "environmentVariables":"a=b",
  "health":100,
  "project":null,
  "interval": 300
}
```

## Get All Tests

```shell
curl "http://syntinel/api/v1/test/all"
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "http://192.168.1.2/api/v1/test/all");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var tests = JSON.parse(xmlHttp.responseText);
    console.log(tests);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id":1,
      "name":"The greatest song in the world",
      "script":"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover",
      "dockerfile":"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh",
      "environmentVariables":"a=b",
      "health":100,
      "project":null,
      "interval": null
    }
  ]
}
```

This endpoint retrieves all tests.

### HTTP Request

`GET http://syntinel/api/v1/test/all`

### Query Parameters
@TODO

## Get All Tests For a Given Project

```shell
curl "http://syntinel/api/v1/test/all?project=1"
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "http://192.168.1.2/api/v1/test/all?project=1");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var tests = JSON.parse(xmlHttp.responseText);
    console.log(tests);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id":1,
      "name":"The greatest song in the world",
      "script":"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover",
      "dockerfile":"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh",
      "environmentVariables":"a=b",
      "health":100,
      "project":"The World's Best Software Project",
      "interval": null
    }
  ]
}
```

This endpoint retrieves all tests.

### HTTP Request

`GET http://syntinel/api/v1/test/all`

### Query Parameters
test integer

## Get a Specific Test

```shell
curl "http://syntinel/api/v1/test/1"
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "http://192.168.1.2/api/v1/test/1");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var test = JSON.parse(xmlHttp.responseText);
    console.log(test);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "id":1,
  "name":"The greatest song in the world",
  "script":"#!/usr/bin/env bash\ngit clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover",
  "dockerfile":"FROM docker.io/centos\n\nMAINTAINER Christopher Henderson\n\nRUN yum install -y go git wget\nCOPY script.sh $HOME/script.sh\nCMD chmod +x script.sh && ./script.sh",
  "environmentVariables":"a=b",
  "health":100,
  "project":null,
  "interval": 300
}
```

This endpoint retrieves a specific test.


### HTTP Request

`GET http://syntinel/api/v1/test/<ID>`

### URL Parameters

Parameter | Description
--------- | -----------
ID | The ID of the test to retrieve

## Schedule a Test

Scheduling a test means that the "interval" field of a Test is some non-null integer value that represents the interval, in seconds, that the test will run. This can be achieved either with a POST or a PATCH.

```shell
curl -H "Content-Type:application/json" -X PATCH http://192.168.1.2/api/v1/test/1 -d '{"interval": 10}'
```

```javascript
var test = {"interval": 10}
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("PATCH", "http://192.168.1.2/api/v1/test/");
xmlHttp.setRequestHeader("Content-Type", "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var test = JSON.parse(xmlHttp.responseText);
    console.log(test);
  }
};
xmlHttp.send(JSON.stringify(test));
```

### HTTP Request

`GET http://syntinel/api/v1/test/<ID>`

### URL Parameters

Parameter | Description
--------- | -----------
ID | The ID of the test to retrieve

## Unschedule a Test

Unscheduling a test means that the "interval" field of a Test is null. This can be achieved either with a POST or a PATCH.

```shell
curl -H "Content-Type:application/json" -X PATCH http://192.168.1.2/api/v1/test/1 -d '{"interval": null}'
```

```javascript
var test = {"interval": null}
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("PATCH", "http://192.168.1.2/api/v1/test/");
xmlHttp.setRequestHeader("Content-Type", "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var test = JSON.parse(xmlHttp.responseText);
    console.log(test);
  }
};
xmlHttp.send(JSON.stringify(test));
```

### HTTP Request

`GET http://syntinel/api/v1/test/<ID>`

### URL Parameters

Parameter | Description
--------- | -----------
ID | The ID of the test to retrieve

# Projects

## Create a Project

Creating a project requires JSON body formatted as follows:

name: String naming the project.<br>

```shell
curl -X POST http://192.168.1.2/api/v1/project/ -d "name=UltimateCode"
```

```javascript
var test = {
  "name":"UltimateCode",
}
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("POST", "http://192.168.1.2/api/v1/project/");
xmlHttp.setRequestHeader("Content-Type", "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var project = JSON.parse(xmlHttp.responseText);
    console.log(project);
  }
};
xmlHttp.send(JSON.stringify(project));
```

> The above command returns JSON structured like this:

```json
{
  "id": 1,
  "tests": [],
  "name":"UltimateCode"
}
```

## Get All Projects

```shell
curl "http://syntinel/api/v1/project/all"
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "http://192.168.1.2/api/v1/project/all");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var projects = JSON.parse(xmlHttp.responseText);
    console.log(projects);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": 1,
      "tests": [1],
      "name": "UltimateCode"
    }
  ]
}  
```

This endpoint retrieves all projects.

### HTTP Request

`GET http://syntinel/api/v1/project/all`

## Get a Specific Project

```shell
curl "http://syntinel/api/v1/project/1"
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "http://192.168.1.2/api/v1/project/1");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var project = JSON.parse(xmlHttp.responseText);
    console.log(project);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "id": 1,
  "tests": [1, 24, 89],
  "name": "UltimateCode"}
```

This endpoint retrieves a specific project.


### HTTP Request

`GET http://syntinel/api/v1/project/<ID>`

# Test Runs

## Create a Test Run

```shell
curl -X POST -H 'Content-Type:application/json' http://192.168.1.2/api/v1/testrun/ \
-d '{"test": 1}'
```

```javascript
var data = {
  "test": 1
}
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("POST", "http://192.168.1.2/api/v1/testrun/", false ); // false for synchronous request
xmlHttp.setRequestHeader('Content-Type', "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var testRun = JSON.parse(xmlHttp.responseText);
    console.log(xmlHttp.responseText);
    console.log("Test Run ID is: " + testRun["id"]);
  }
};
xmlHttp.send(JSON.stringify(data));
```

> The above command returns JSON structured like this:

```json
{
  "id":71,
  "log":"",
  "successful":null,
  "test":1
}
```

## Get All Test Runs

```shell
curl -X POST -H 'Content-Type:application/json' http://192.168.1.2/api/v1/testrun/all
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("POST", "http://192.168.1.2/api/v1/testrun/all", false ); // false for synchronous request
xmlHttp.setRequestHeader('Content-Type', "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var testRun = JSON.parse(xmlHttp.responseText);
    console.log(xmlHttp.responseText);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "count":1,
  "next":null,
  "previous":null,
  "results":[
    {
      "id":1,
      "log":"Sending build context to Docker daemon 3.072 kB\r\nStep 1 : FROM docker.io/centos\n ---> 67591570dd29\nStep 2 : MAINTAINER Christopher Henderson\n ---> Using cache\n ---> d33126dcfbfc\nStep 3 : RUN yum install -y go git wget\nCloning into 'TestTheTester'...\n ---> Using cache\n ---> b4da48a3d37c\nStep 4 : COPY script.sh $HOME/script.sh\n ---> Using cache\n ---> 40743fa214db\nStep 5 : CMD chmod +x script.sh && ./script.sh\n ---> Using cache\n ---> 1da95f92dcd2\nSuccessfully built 1da95f92dcd2\n=== RUN   TestPow\n--- PASS: TestPow (0.00s)\n=== RUN   TestSquare\n--- PASS: TestSquare (0.00s)\nPASS\ncoverage: 46.2% of statements\n",
      "error":"",
      "status":null,
      "successful":true,
      "test":1
    }
  ]
}
```

## Get All Test Runs For a Given Test

```shell
curl -X POST -H 'Content-Type:application/json' http://192.168.1.2/api/v1/testrun/all?test=1
```

```javascript
var xmlHttp = new XMLHttpRequest();
xmlHttp.open("POST", "http://192.168.1.2/api/v1/testrun/all?test=1", false ); // false for synchronous request
xmlHttp.setRequestHeader('Content-Type', "application/json");
xmlHttp.onreadystatechange = function() {
  if (xmlHttp.readyState === 4)  {
    var testRun = JSON.parse(xmlHttp.responseText);
    console.log(xmlHttp.responseText);
  }
};
xmlHttp.send(null);
```

> The above command returns JSON structured like this:

```json
{
  "count":1,
  "next":null,
  "previous":null,
  "results":[
    {
      "id":1,
      "log":"Sending build context to Docker daemon 3.072 kB\r\nStep 1 : FROM docker.io/centos\n ---> 67591570dd29\nStep 2 : MAINTAINER Christopher Henderson\n ---> Using cache\n ---> d33126dcfbfc\nStep 3 : RUN yum install -y go git wget\nCloning into 'TestTheTester'...\n ---> Using cache\n ---> b4da48a3d37c\nStep 4 : COPY script.sh $HOME/script.sh\n ---> Using cache\n ---> 40743fa214db\nStep 5 : CMD chmod +x script.sh && ./script.sh\n ---> Using cache\n ---> 1da95f92dcd2\nSuccessfully built 1da95f92dcd2\n=== RUN   TestPow\n--- PASS: TestPow (0.00s)\n=== RUN   TestSquare\n--- PASS: TestSquare (0.00s)\nPASS\ncoverage: 46.2% of statements\n",
      "error":"",
      "status":null,
      "successful":true,
      "test":1
    }
  ]
}
```

# Executor

## Register an Executor

Registering an executor requires JSON body formatted as follows:

hostName: String<br>
port: String<br>
Scheme: String. Optional. Default is "http"<br>

```shell
curl -X POST -H "Content-Type:application/json" http://192.168.1.2/api/v1/executor/ -d '{"port": "9090", "hostName": "localhost"}'
```

> The above command returns JSON structured like this:

```json
{
  "id":1,
  "hostName":"localhost",
  "port":"9090",
  "Scheme":"http"
}
```

## Get All Executors

```shell
curl -X GET http://192.168.1.2/api/v1/executor/all/
```

> The above command returns JSON structured like this:

```json
{
  "count": 2,
  "next": null,
  "previous": null,
  "results": [
    {
      "id":1,
      "hostName":"localhost",
      "port":"9090",
      "Scheme":"http"
    },
    {
      "id":2,
      "hostName":"localhost",
      "port":"9091",
      "Scheme":"http"
    }
  ]
}
```

## Get a Specific Executor

```shell
curl -X GET http://192.168.1.2/api/v1/executor/1
```

> The above command returns JSON structured like this:

```json
{
  "id":1,
  "hostName":"localhost",
  "port":"9090",
  "Scheme":"http"
}
```
