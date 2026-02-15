package com.longhe.learn.mymall.core.util;

import com.longhe.learn.mymall.payment.dao.bo.Business;
import com.longhe.learn.mymall.payment.mapper.generator.po.BusinessPo;
import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;

import static com.longhe.learn.mymall.core.util.Common.cloneObj;
import static org.junit.jupiter.api.Assertions.*;

class CommonTest {

    @Test
    void cloneObj1() {
        LocalDateTime now = LocalDateTime.now();
        BusinessPo po = new BusinessPo();
        po.setId(Long.valueOf(1));
        po.setName("测试");
        po.setCreatorId(Long.valueOf(2));
        po.setCreatorName("管理员");
        po.setGmtCreate(now);
        po.setGmtModified(now.minusDays(Long.valueOf(1)));

        Business bo = cloneObj(po, Business.class);
        assertEquals(Long.valueOf(1), bo.getId());
        assertEquals("测试", bo.getName());
        assertEquals(Long.valueOf(2), bo.getCreatorId());
        assertEquals("管理员",bo.getCreatorName());
        assertEquals(now, bo.getGmtCreate());
        assertEquals(now.minusDays(Long.valueOf(1)), bo.getGmtModified());
    }
}