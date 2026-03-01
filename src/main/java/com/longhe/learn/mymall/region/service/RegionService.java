package com.longhe.learn.mymall.region.service;

import com.longhe.learn.mymall.core.mapper.RedisUtil;
import com.longhe.learn.mymall.core.model.dto.UserDto;
import com.longhe.learn.mymall.region.dao.RegionDao;
import com.longhe.learn.mymall.region.dao.bo.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Transactional
public class RegionService {

    private final Logger logger = LoggerFactory.getLogger(RegionService.class);

    @Value("${oomall.region.timeout}")
    private int timeout;

    private final RegionDao regionDao;
    private final RedisUtil redisUtil;

    @Autowired
    public RegionService(RegionDao regionDao, RedisUtil redisUtil) {
        this.regionDao = regionDao;
        this.redisUtil = redisUtil;
    }

    /**
     * 创建新的子地区
     *
     * @param id region id
     * @param region 下级region对象
     * @param user 登录用户
     * @return 新region对象，带id
     */
    public Region createSubRegions(Long id, Region region, UserDto user) {
        Region parent = this.regionDao.findById(id);
        return parent.createSubRegion(region, user);
    }

}
