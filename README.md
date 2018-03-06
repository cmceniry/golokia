# golokia

Simple jolokia JMX/HTTP wrapper for Go

godocs are not even close to being proper yet.

## How to use

```
domains, err := ListDomains(jolokiaUrl)

beans, err := ListBeans(jolokiaUrl, domainName)

props, err := ListProperties(jolokiaUrl, domainName, propertyName)

val, err := GetAttr(jolokiaUrl, domainName, propertyName, attributeName)

operations, err := ListOperartions(jolokiaUrl, domainName, beanName)

val, err := ExecOp(jolokiaUrl, domainName, operationName, attributesValues)
```
