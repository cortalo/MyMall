package com.longhe.learn.mymall.region.service;

import com.longhe.learn.mymall.core.mapper.RedisUtil;
import com.longhe.learn.mymall.core.model.dto.UserDto;
import com.longhe.learn.mymall.region.dao.RegionDao;
import com.longhe.learn.mymall.region.dao.bo.Region;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.transaction.annotation.Transactional;

import static org.junit.jupiter.api.Assertions.*;


@SpringBootTest
@Transactional
class RegionServiceTest {

    private final static Logger logger =  LoggerFactory.getLogger(RegionServiceTest.class);

    @Autowired
    private RegionService regionService;
    @Autowired
    private RegionDao regionDao;
    @MockBean
    private RedisUtil redisUtil;

    @Test
    void createSubRegions() {
        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        Region region = new Region();
        region.setName("Beijing");

        UserDto userDto = new UserDto();
        userDto.setId(0L);
        userDto.setName("0");
        userDto.setUserLevel(0);
        userDto.setDepartId(0L);

        Region subRegions = regionService.createSubRegions(0L, region, userDto);
        assertNotNull(subRegions);
        assertNotNull(subRegions.getId());
        assertTrue(subRegions.getId() > 0);
        assertEquals(0L, subRegions.getPid());
        assertEquals((byte) 2, subRegions.getLevel());
        assertEquals("Beijing", subRegions.getName());

        Region regionFromDao = regionDao.findById(subRegions.getId());
        assertNotNull(regionFromDao);
        assertEquals(subRegions.getId(), regionFromDao.getId());
        assertEquals(0L, regionFromDao.getPid());
        assertEquals((byte) 2, regionFromDao.getLevel());
        assertEquals("Beijing", regionFromDao.getName());
    }
}