package com.longhe.learn.javaee.core.aop;

import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Pointcut;

@Aspect
public class CommonPointCuts {

    @Pointcut("execution(public com.longhe.learn.javaee.core.model.ReturnObject com.longhe.learn..controller..*.*(..))")
    public void controllers() {
    }

    @Pointcut("execution(public * com.longhe.learn..dao..*.*(..))")
    public void daos(){
    }

    @Pointcut("@annotation(com.longhe.learn.javaee.core.aop.Audit)")
    public void auditAnnotation() {
    }
}
