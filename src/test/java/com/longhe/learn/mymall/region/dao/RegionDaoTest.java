package com.longhe.learn.mymall.region.dao;

import com.longhe.learn.mymall.core.mapper.RedisUtil;
import com.longhe.learn.mymall.region.dao.bo.Region;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
class RegionDaoTest {

    private final static Logger logger = LoggerFactory.getLogger(RegionDaoTest.class);

    @Autowired
    private RegionDao regionDao;
    @MockBean
    private RedisUtil redisUtil;

    @Test
    void findById() {
        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        Region region = regionDao.findById(0L);
        assertNotNull(region);
        assertEquals(0L, region.getId());
        assertEquals(-1L, region.getPid());
        assertEquals("China", region.getName());
    }

    @Test
    void findById1() {
        Region region = new Region();
        region.setId(777L);

        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(region);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        Region bo = regionDao.findById(0L);
        assertNotNull(bo);
        assertEquals(777L, bo.getId());
        assertNull(bo.getName());
    }
}