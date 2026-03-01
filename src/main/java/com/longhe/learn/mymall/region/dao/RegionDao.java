package com.longhe.learn.mymall.region.dao;

import com.longhe.learn.mymall.core.exception.BusinessException;
import com.longhe.learn.mymall.core.mapper.RedisUtil;
import com.longhe.learn.mymall.core.model.ReturnNo;
import com.longhe.learn.mymall.core.model.dto.UserDto;
import com.longhe.learn.mymall.core.util.Common;
import com.longhe.learn.mymall.region.dao.bo.Region;
import com.longhe.learn.mymall.region.mapper.RegionPoMapper;
import com.longhe.learn.mymall.region.mapper.po.RegionPo;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Repository;

import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.Optional;

import static com.longhe.learn.mymall.core.util.Common.cloneObj;

@Repository
public class RegionDao {

    private final static Logger logger = LoggerFactory.getLogger(RegionDao.class);
    private final static String KEY = "R%d";

    private final static String PARENT_KEY = "RP%d";

    @Value("${oomall.region.timeout}")
    private int timeout;

    private final RegionPoMapper regionPoMapper;
    private final RedisUtil redisUtil;

    @Autowired
    public RegionDao(RegionPoMapper regionPoMapper, RedisUtil redisUtil) {
        this.regionPoMapper = regionPoMapper;
        this.redisUtil = redisUtil;
    }

    public void build(Region bo) {
        bo.setRegionDao(this);
    }

    public Region build(RegionPo po, Optional<String> redisKey) {
        Region bo = cloneObj(po, Region.class);
        this.build(bo);
        redisKey.ifPresent(key -> redisUtil.set(key, bo, timeout));
        return bo;
    }

    public Region findById(Long id) {
        logger.debug("findById: id = {}", id);
        if (null == id) {
            throw new IllegalArgumentException("id can not be null");
        }
        String key = String.format(KEY, id);
        Region bo = (Region) redisUtil.get(key);
        if (null != bo) {
            logger.debug("findById: hit in redis key = {}, region = {}", key, bo);
            this.build(bo);
            return bo;
        } else {
            Optional<RegionPo> ret = regionPoMapper.findById(id);
            if (ret.isPresent()) {
                logger.debug("findById: retrieve from database region = {}", ret.get());
                return this.build(ret.get(), Optional.of(key));
            } else {
                throw new BusinessException(ReturnNo.RESOURCE_ID_NOTEXIST, String.format(ReturnNo.RESOURCE_ID_NOTEXIST.getMessage(), "地区", id));
            }
        }
    }

    public Region insert(Region bo, UserDto user) {
        bo.setId(null);
        bo.setCreator(user);
        bo.setGmtCreate(LocalDateTime.now());
        RegionPo po = cloneObj(bo, RegionPo.class);
        logger.debug("save: po = {}", po);
        po = regionPoMapper.save(po);
        bo.setId(po.getId());
        return bo;
    }
}
