# Desired state
A single application which can be used to serialize user specified payloads
using specified serializers (in our case AVRO or JSON) and publish that message 
dynamically to the topic specified by the user. Additionally, would like the topic
to be allowed to be usable for multiple type publications.


## Setup
Below values can be changed in the `application.yaml`
- Port `8089` free to run on
- Kafka cluster running on `127.0.0.1:9092`
- Schema Registry running `127.0.0.1:8081`
  - With some test schemas preloaded
- Postman (or equivalent tool) to POST data to app

## Request
    POST http://127.0.0.1:8089/(avro|json)?
    {
        "topic": "<your topic here>",
        "avroSource": "<your schema subject>",
        "payload": {
            "<examplePayloadKey1>": "<examplePayloadValue1>",
            "<examplePayloadKey2>": "<examplePayloadValue2>"
        },
        "headers": {
            "<exampleHeaderKey1>": "<exampleHeaderValue1>",
            "<exampleHeaderKey2>": "<exampleHeaderValue2>"
        }
    }

## Currently
Using either `AVRO` or `JSON` as the spring profile produces the desired result: namely that the event is published
to kafka using the correct serializer.

However, when the spring profile is set to `BOTH` then we get the below error when trying to publish:
```
2022-05-13 15:40:42.332 ERROR 9997 --- [nio-8089-exec-2] o.a.c.c.C.[.[.[/].[dispatcherServlet]    : Servlet.service() for servlet [dispatcherServlet] in context with path [] threw exception [Request processing failed; nested exception is java.lang.IllegalStateException: A default binder has been requested, but there is more than one binder available for 'org.springframework.messaging.MessageChannel' : json,avro, and no default binder has been set.] with root cause

java.lang.IllegalStateException: A default binder has been requested, but there is more than one binder available for 'org.springframework.messaging.MessageChannel' : json,avro, and no default binder has been set.
	at org.springframework.cloud.stream.binder.DefaultBinderFactory.doGetBinder(DefaultBinderFactory.java:210) ~[spring-cloud-stream-3.2.3.jar:3.2.3]
	at org.springframework.cloud.stream.binder.DefaultBinderFactory.getBinder(DefaultBinderFactory.java:151) ~[spring-cloud-stream-3.2.3.jar:3.2.3]
	at org.springframework.cloud.stream.function.StreamBridge.resolveBinderTargetType(StreamBridge.java:305) ~[spring-cloud-stream-3.2.3.jar:3.2.3]
	at org.springframework.cloud.stream.function.StreamBridge.send(StreamBridge.java:221) ~[spring-cloud-stream-3.2.3.jar:3.2.3]
	at com.example.MessageUtils.sendMessage(MessageUtils.java:39) ~[classes/:na]
	at com.example.controllers.ControllerBoth.postJSON(ControllerBoth.java:28) ~[classes/:na]
	at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(DirectMethodHandleAccessor.java:104) ~[na:na]
	at java.base/java.lang.reflect.Method.invoke(Method.java:577) ~[na:na]
	at org.springframework.web.method.support.InvocableHandlerMethod.doInvoke(InvocableHandlerMethod.java:205) ~[spring-web-5.3.18.jar:5.3.18]
	at org.springframework.web.method.support.InvocableHandlerMethod.invokeForRequest(InvocableHandlerMethod.java:150) ~[spring-web-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.mvc.method.annotation.ServletInvocableHandlerMethod.invokeAndHandle(ServletInvocableHandlerMethod.java:117) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.invokeHandlerMethod(RequestMappingHandlerAdapter.java:895) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.handleInternal(RequestMappingHandlerAdapter.java:808) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.mvc.method.AbstractHandlerMethodAdapter.handle(AbstractHandlerMethodAdapter.java:87) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:1067) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:963) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1006) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at org.springframework.web.servlet.FrameworkServlet.doPost(FrameworkServlet.java:909) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at javax.servlet.http.HttpServlet.service(HttpServlet.java:681) ~[tomcat-embed-core-9.0.60.jar:4.0.FR]
	at org.springframework.web.servlet.FrameworkServlet.service(FrameworkServlet.java:883) ~[spring-webmvc-5.3.18.jar:5.3.18]
	at javax.servlet.http.HttpServlet.service(HttpServlet.java:764) ~[tomcat-embed-core-9.0.60.jar:4.0.FR]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:227) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53) ~[tomcat-embed-websocket-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100) ~[spring-web-5.3.18.jar:5.3.18]
	at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117) ~[spring-web-5.3.18.jar:5.3.18]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93) ~[spring-web-5.3.18.jar:5.3.18]
	at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117) ~[spring-web-5.3.18.jar:5.3.18]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.springframework.boot.actuate.metrics.web.servlet.WebMvcMetricsFilter.doFilterInternal(WebMvcMetricsFilter.java:96) ~[spring-boot-actuator-2.6.6.jar:2.6.6]
	at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117) ~[spring-web-5.3.18.jar:5.3.18]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:201) ~[spring-web-5.3.18.jar:5.3.18]
	at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:117) ~[spring-web-5.3.18.jar:5.3.18]
	at org.apache.catalina.core.ApplicationFilterChain.internalDoFilter(ApplicationFilterChain.java:189) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:162) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:197) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:97) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:541) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:135) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:92) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:78) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:360) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:399) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:65) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:889) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1743) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:49) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1191) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:659) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:61) ~[tomcat-embed-core-9.0.60.jar:9.0.60]
	at java.base/java.lang.Thread.run(Thread.java:833) ~[na:na]
```

### Partial really BAD work around is to use AOP
Set the spring profile to `BOTH,ASPECT`

I created an aspect to try and prevent the "StreamBridge.resolveBinderTargetType" method from
throwing the IllegalStateException, which does allow the application to get past the exception
and allows the application to publish as expected to the topic with the given caveat:  

    Once a Topic has been used for AVRO or JSON it can no longer be used for the other without an application restart
    - if the topic was used first for AVRO and you attempt to publish JSON through it then you get `org.apache.kafka.common.errors.SerializationException: In configuration value.subject.name.strategy = io.confluent.kafka.serializers.subject.RecordNameStrategy, the message value must only be a record schema`
    - if the topic was used first for JSON and you attempt to publish AVRO through it then you get `java.lang.ClassCastException: class org.apache.avro.generic.GenericData$Record cannot be cast to class [B (org.apache.avro.generic.GenericData$Record is in unnamed module of loader 'app'; [B is in module java.base of loader 'bootstrap')`

This functionality is not desirable either as we would like to have this application be 
flexible and not cause a topic to become bound to one serialization type.


## Notes
- The `converter` and `org.apache.avro.io` packages placed at `src/main/java` are used to represent the below dependency
while preventing security vulnerabilities.

```
<!-- https://mvnrepository.com/artifact/tech.allegro.schema.json2avro/converter -->
<dependency>
<groupId>tech.allegro.schema.json2avro</groupId>
<artifactId>converter</artifactId>
<version>0.2.2</version>
</dependency>
```