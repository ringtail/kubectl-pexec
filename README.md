# kubectl-pexec 
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Build Status](https://travis-ci.org/ringtail/kubectl-pexec.svg?branch=master)](https://travis-ci.org/ringtail/pexec)     

`kubectl-pexec` is inspired by `pssh`. When you want to exec some commands in several pods of a same Deployment. It's very diffcult to do it in kubernetes with kubectl. But it is very common for ops-man🔧👱. 

## Usage 
```sh
# assumes you have a working KUBECONFIG
$ GO111MODULE="on" go build cmd/kubectl-pexec.go
# place the built binary somewhere in your PATH
$ cp ./kubectl-pexec /usr/local/bin

# you can now begin using this plugin as a regular kubectl command:
```

#### show uname of all pods of a deploy 
```sh 
$ kubectl pexec deploy nginx-deployment-basic -n default "uname -a" 

[nginx-deployment-basic-64fc4c755d-7h49v] Linux nginx-deployment-basic-64fc4c755d-7h49v 4.19.57-15.1.al7.x86_64 #1 SMP Thu Aug 29 13:46:41 CST 2019 x86_64 GNU/Linux
[nginx-deployment-basic-64fc4c755d-zqqv8] Linux nginx-deployment-basic-64fc4c755d-zqqv8 4.19.57-15.1.al7.x86_64 #1 SMP Thu Aug 29 13:46:41 CST 2019 x86_64 GNU/Linux
All pods execution done in 0.547s
```

#### check the nginx config of all pods 
```sh 
$ kubectl pexec deploy nginx-deployment-basic -n default cat /etc/nginx/nginx.conf

[nginx-deployment-basic-64fc4c755d-7h49v]
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;
}
[nginx-deployment-basic-64fc4c755d-zqqv8]
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;
}
All pods execution done in 0.547s
```

#### debug or tuning in pod 
```sh 
$ kubectl pexec deploy nginx-deployment-basic "netstat -apn"

[nginx-deployment-basic-64fc4c755d-7h49v]Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN      1/nginx: master pro
Active UNIX domain sockets (servers and established)
Proto RefCnt Flags       Type       State         I-Node   PID/Program name     Path
unix  3      [ ]         STREAM     CONNECTED     243464938 1/nginx: master pro
unix  3      [ ]         STREAM     CONNECTED     243464939 1/nginx: master pro
[nginx-deployment-basic-64fc4c755d-zqqv8]Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN      1/nginx: master pro
Active UNIX domain sockets (servers and established)
Proto RefCnt Flags       Type       State         I-Node   PID/Program name     Path
unix  3      [ ]         STREAM     CONNECTED     251539135 1/nginx: master pro
unix  3      [ ]         STREAM     CONNECTED     251539136 1/nginx: master pro
All pods execution done in 0.547s
```

## Use Cases 
* Batch commands execution 
* Problem diagnose 
* Performance tuning 
* Configuration alteration 

## Frequently asked question
* Why my command can not be executed?   
You can try to wrap the command with a colon. some special characters may break the command parser.
                                             
* Why the stdout is mixed and out-of-order?
`kubectl-pexec` bond the stdout to multi remote stdout. So if the stream is continuous. The stdout may be out-of-order.

## License
This software is released under the Apache 2.0 license.
