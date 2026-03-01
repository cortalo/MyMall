package com.longhe.learn.mymall.core.aop;

import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Pointcut;

@Aspect
public class CommonPointCuts {

    @Pointcut("execution(public com.longhe.learn.mymall.core.model.ReturnObject com.longhe.learn.mymall..controller..*.*(..))")
    public void controllers() {
    }
}
