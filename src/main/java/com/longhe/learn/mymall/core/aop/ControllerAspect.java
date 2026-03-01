package com.longhe.learn.mymall.core.aop;


import com.longhe.learn.mymall.core.exception.BusinessException;
import com.longhe.learn.mymall.core.model.ReturnNo;
import com.longhe.learn.mymall.core.model.ReturnObject;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.Signature;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

@Aspect
@Component
@Order(10)
public class ControllerAspect {

    private final Logger logger = LoggerFactory.getLogger(ControllerAspect.class);

    @Around("com.longhe.learn.mymall.core.aop.CommonPointCuts.controllers()")
    public Object doAround(ProceedingJoinPoint jp) throws Throwable {
        ReturnObject retVal = null;

        MethodSignature ms = (MethodSignature) jp.getSignature();
        HttpServletRequest request = ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes()).getRequest();
        HttpServletResponse response = ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes()).getResponse();

        String[] paramNames = ms.getParameterNames();
        logger.debug("doAround: method = {}, paramNames = {}", ms.getName(), paramNames);
        Object[] args = jp.getArgs();
        try {
            retVal = (ReturnObject) jp.proceed(args);
        } catch (BusinessException exception) {
            logger.info("doAround: BusinessException, errno = {}", exception.getErrno());
            retVal = new ReturnObject(exception.getErrno(), exception.getMessage());
        }

        ReturnNo code = retVal.getCode();
        logger.debug("doAround: jp = {}, code = {}", jp.getSignature().getName(), code);
        changeHttpStatus(code, response);

        return retVal;
    }


    /**
     * 根据code，修改reponse的HTTP Status code
     *
     * @param code
     * @param response
     */
    private void changeHttpStatus(ReturnNo code, HttpServletResponse response) {
        switch (code) {
            case CREATED:
                //201:
                response.setStatus(HttpServletResponse.SC_CREATED);
                break;

            case RESOURCE_ID_NOTEXIST:
                // 404：资源不存在
                response.setStatus(HttpServletResponse.SC_NOT_FOUND);
                break;

            case AUTH_INVALID_JWT:
            case AUTH_JWT_EXPIRED:
            case AUTH_NEED_LOGIN:
            case AUTH_ID_NOTEXIST:
            case AUTH_USER_FORBIDDEN:
            case AUTH_INVALID_ACCOUNT:
                // 401
                response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                break;

            case INTERNAL_SERVER_ERR:
                // 500：数据库或其他严重错误
                response.setStatus(HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
                break;

            case FIELD_NOTVALID:
            case IMG_FORMAT_ERROR:
            case IMG_SIZE_EXCEED:
            case INCONSISTENT_DATA:
                // 400
                response.setStatus(HttpServletResponse.SC_BAD_REQUEST);
                break;

            case RESOURCE_ID_OUTSCOPE:
            case FILE_NO_WRITE_PERMISSION:
            case AUTH_NO_RIGHT:
                // 403
                response.setStatus(HttpServletResponse.SC_FORBIDDEN);
                break;
            default:
                response.setStatus(HttpServletResponse.SC_OK);
        }
        response.setContentType("application/json;charset=UTF-8");
    }

}
