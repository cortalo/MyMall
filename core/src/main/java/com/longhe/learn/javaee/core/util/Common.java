package com.longhe.learn.javaee.core.util;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

/**
 * 通用工具类
 *
 * @author Ming Qiu
 **/
public class Common {
    private static Logger logger = LoggerFactory.getLogger(Common.class);

    private static Pattern setPattern = Pattern.compile("set[A-Z][a-z][a-zA-Z0-9]*");

    private static Pattern getPattern = Pattern.compile("get[A-Z][a-z][a-zA-Z0-9]*");

    private static Pattern firstCharPattern = Pattern.compile("^.");

    /**
     * 生成九位数序号
     * 要保证同一服务的不同实例生成出的序号是不同的
     * @param  platform 机器号 如果一个服务有多个实例，机器号需不同，目前从1至36
     * @return 序号
     */
    public static String genSeqNum(int platform) {
        int maxNum = 36;
        int i;

        DateTimeFormatter dtf = DateTimeFormatter.ofPattern("yyyyMMddHHmmssS");
        LocalDateTime localDateTime = LocalDateTime.now();
        String strDate = localDateTime.format(dtf);
        StringBuffer sb = new StringBuffer(strDate);

        int count = 0;
        char[] str = {'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K',
                'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
                'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'};
        Random r = new Random();
        while (count < 2) {
            i = Math.abs(r.nextInt(maxNum));
            if (i >= 0 && i < str.length) {
                sb.append(str[i]);
                count++;
            }
        }
        if (platform > 36){
            platform = 36;
        } else if (platform < 1){
            platform = 1;
        }

        sb.append(str[platform-1]);
        return sb.toString();
    }


}