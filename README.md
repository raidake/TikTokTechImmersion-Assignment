
# Tiktok Tech Immersion 2023 - IM Service

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

Demo and template retrieved from [TikTok Tech Immersion 2023 Backend Assignment](https://github.com/TikTokTechImmersion/assignment_demo_2023). 
The goal of the assignment is to develop an IM service that is implemented with a set of specific APIs using Golang. 

For my assignment, the messages are stored in a MySQL Database. The HTTP service is connected to the RPC service, which is then connected to the MySQL Database (in a separate container).

## Tools used

These are the tools that I used to complete the assignment
* [Docker Desktop](https://www.docker.com/products/docker-desktop)
*  [Go](https://golang.org/doc/install)
* [JMeter](https://jmeter.apache.org/download_jmeter.cgi)
* [Postman](https://www.postman.com/downloads/)


## How to Start

Run the command below to start up the containers
```
docker-compose up -d
```
\
Sending messages
*Example: Sending "hi" from a to b*
```
curl -X POST 'localhost:8080/api/send?sender=a&receiver=b&text=hi'
```
\
Retrieving chat logs
*Example: Retrieving the messages between a and b*
```
curl 'localhost:8080/api/pull?chat=a%3Ab'
```
*By default, cursor = 0, limit = 10 and reverse = false*
To specify, you can run the command like the example below,
```
curl 'localhost:8080/api/pull?chat=a%3Ab&cursor=0&limit=4&reverse=true'
```
where cursor is set to 0, limit is set 4 and reverse is set to true for ascending order based on timestamp.

# Stress Testing Results
I used JMeter to perform concurrent testing to see how well the service will run based on the number of users

1 User:
![Test for 1 User](/img/1user.jpg "Test for 1 User")

20 User:
![Test for 20 Users](/img/20user.jpg "Test for 20 Users")

100 User:
![Test for 100 Users](/img/100user.jpg "Test for 100 Users")