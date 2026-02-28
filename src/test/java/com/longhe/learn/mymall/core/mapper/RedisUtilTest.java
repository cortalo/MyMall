package com.longhe.learn.mymall.core.mapper;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
class RedisUtilTest {

    @Autowired
    private RedisUtil redisUtil;

    @Test
    void testConnection() {
        String key = "test:connection";
        String value = "hello_redis";

        // 写入
        redisUtil.set(key, value, 60);

        // 读取
        Object result = redisUtil.get(key);

        System.out.println("Redis连接成功，读取到的值：" + result);
        assertEquals(value, result);

        // 清理
        redisUtil.del(key);
    }

}