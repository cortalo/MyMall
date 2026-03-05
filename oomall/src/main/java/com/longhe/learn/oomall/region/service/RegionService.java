package com.longhe.learn.oomall.region.service;

import com.longhe.learn.javaee.core.mapper.RedisUtil;
import com.longhe.learn.javaee.core.model.UserToken;
import com.longhe.learn.oomall.region.dao.RegionDao;
import com.longhe.learn.oomall.region.dao.bo.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

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
    public Region createSubRegions(Long id, Region region, UserToken user) {
        Region parent = this.regionDao.findById(id);
        return parent.createSubRegion(region, user);
    }

    /**
     * 通过id查找地区
     *
     * @param id 地区id
     * @return RegionDto
     */
    public Region findById(Long id) {
        logger.debug("findRegionById: id = {}", id);
        return this.regionDao.findById(id);
    }

    public List<Region> retrieveParentsRegionsById(Long id) {
        logger.debug("retrieveParentsRegionsById: id = {}", id);
        Region region = this.regionDao.findById(id);
        return region.getAncestors();
    }

    /**
     * abandon region by id, will abandon subRegions
     * @param id region id
     * @param user userToken
     */
    public void deleteRegion(Long id, UserToken user) {
        Region region = this.regionDao.findById(id);
        List<String> keys = region.abandon(user);
        this.redisUtil.del(keys.toArray(new String[0]));
    }
}
