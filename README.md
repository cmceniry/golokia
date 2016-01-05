# golokia

```
< golokia >
 --------
   \
    \   
       _____  -----------------  _____
      |     \/  ----     ----  \/     |
       \  ==/  /    |   |    \  \==  /
        \--/  | ()  |   |()   |  \--/
          /    \---/ ___ \---/    \
         |         _(===)_         |
        |         (__/--\__)        |
         |          |_||_|         |
          \                       /   
```

Golokia is a Simple jolokia JMX/HTTP wrapper for Go. It supports jolokia proxy setup as well as direct host connection.

[Golokia GoDoc] (https://godoc.org/github.com/swarvanusg/golokia)

### Version
0.1.0

### Usage

#### Step 1 : Get It
To get the GoPlug install Go and execute the below command 
```
go get github.com/swarvanusg/GoPlug
```

#### Step 2 : Initiate a client
```go
client := NewJolokiaClient("http://" + proxyhost + ":" + proxyport + "/" + jolokia)

client.SetTarget(targetHost + ":" + targetPort)
```

#### Step 3 : Use the client for Getting Info
```go
beans, err := client.ListBeans("java.lang")

props, err := client.ListProperties("java.lang", []string{"type=Threading"})

val, err := client.GetAttr("java.lang", []string{"type=Threading"}, "PeakThreadCount")
```

#### Step 4 : Use it in a direct approach
```go
domains, err := ListDomains(jolokiaUrl)

beans, err := ListBeans(jolokiaUrl, domainName)

props, err := ListProperties(jolokiaUrl, domainName, propertyName)

val, err := GetAttr(jolokiaUrl, domainName, propertyName, attributeName)
```

#### Current Status:
The Golokia build is success
The test cases are passing 

