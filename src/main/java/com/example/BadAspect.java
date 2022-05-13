package com.example;

import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.cloud.stream.binder.Binder;
import org.springframework.cloud.stream.binder.Binding;
import org.springframework.cloud.stream.binder.ConsumerProperties;
import org.springframework.cloud.stream.binder.ProducerProperties;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Service;

@Service
@Aspect
@Profile("ASPECT")
public class BadAspect {

    @Around("execution(* org.springframework.cloud.stream.binder.DefaultBinderFactory.getBinder(..))")
    public <T, C extends ConsumerProperties, P extends ProducerProperties> Object stubGetBinder(ProceedingJoinPoint proceedingJoinPoint) {
        Object[] args = proceedingJoinPoint.getArgs();
        try {
            return proceedingJoinPoint.proceed(args);
        } catch (Throwable t) {
            t.printStackTrace();
            return new Binder<T, C, P> () {
                @Override
                public Binding<T> bindConsumer(String name, String group, T inboundBindTarget, C consumerProperties) {
                    return null;
                }

                @Override
                public Binding<T> bindProducer(String name, T outboundBindTarget, P producerProperties) {
                    return null;
                }

            };
        }
    }

    /**
     * This is my desired AOP to solve the issue, but since the StreamBridge class is final I cannot.
     * @param proceedingJoinPoint
     * @return
     * @throws Throwable
     */
//    @Around("execution(* org.springframework.cloud.stream.function.StreamBridge.resolveBinderTargetType(..))")
    public Object stubResolveBinderTargetType(ProceedingJoinPoint proceedingJoinPoint) throws Throwable {
        return "kafka";
    }

}
