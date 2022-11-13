System is written in GO. 
Server will restart if it fail to start first.




PS C:\Users\\Documents\_project_\devopstuni> uname -a; docker --version; docker-compose --version
uname : The term 'uname' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
At line:1 char:1
+ uname -a; docker --version; docker-compose --version
+ ~~~~~
    + CategoryInfo          : ObjectNotFound: (uname:String) [], CommandNotFoundException
    + FullyQualifiedErrorId : CommandNotFoundException

Docker version 20.10.13, build a224086
docker-compose version 1.29.2, build 5becea4c


Not much learning about rabbit message system. Learned how to use GO/Golang.

Messaging system use useful mainly when creating microservice system or service buss system. 
So when there is asynchronous messaging between service  with one to many or many to one relation.
