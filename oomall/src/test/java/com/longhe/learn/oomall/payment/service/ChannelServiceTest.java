package com.longhe.learn.oomall.payment.service;

import com.longhe.learn.javaee.core.mapper.RedisUtil;
import com.longhe.learn.oomall.OomallApplication;
import com.longhe.learn.oomall.payment.dao.bo.Account;
import com.longhe.learn.oomall.region.service.RegionService;
import lombok.extern.slf4j.Slf4j;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest(classes = OomallApplication.class)
@Transactional
@Slf4j
class ChannelServiceTest {

    @Autowired
    private ChannelService channelService;
    @MockBean
    private RedisUtil redisUtil;

    @Test
    void retrieveValidAccount() {
        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);
        List<Account> accounts = channelService.retrieveValidAccount(1001L, 1, 10);
        assertNotNull(accounts);
        assertEquals(1, accounts.size());
        assertEquals(1L, accounts.get(0).getId());
        assertEquals(1001L, accounts.get(0).getShopId());
    }
}