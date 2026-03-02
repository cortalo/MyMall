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
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Repository;

import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

import static com.longhe.learn.mymall.core.model.Constants.IDNOTEXIST;
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

    public List<Region> retrieveParentsRegions(Region region) {
        String key = String.format(PARENT_KEY, region.getId());
        List<Long> parentIds = (List<Long>) redisUtil.get(key);
        if (null != parentIds) {
            return parentIds.stream().map(this::findById).collect(Collectors.toList());
        }
        List<Region> regions = new ArrayList<>();
        while (regions.size() < 10 && !Region.ROOT_PID.equals(region.getPid())) {
            region = this.findById(region.getPid());
            regions.add(region);
        }
        this.redisUtil.set(key, new ArrayList<>(regions.stream().map(Region::getId).collect(Collectors.toList())), timeout);
        return regions;
    }

    public String save(Region bo, UserDto user) {
        bo.setModifier(user);
        bo.setGmtModified(LocalDateTime.now());
        RegionPo po = cloneObj(bo, RegionPo.class);
        RegionPo updatePo = regionPoMapper.save(po);
        if (IDNOTEXIST.equals(updatePo.getId())) {
            throw new BusinessException(ReturnNo.RESOURCE_ID_NOTEXIST, String.format(ReturnNo.RESOURCE_ID_NOTEXIST.getMessage(), "地区", bo.getId()));
        }
        return String.format(KEY, bo.getId());
    }

    public List<Region> retrieveSubRegionsById(Long pid, Boolean all, Integer page, Integer pageSize) {
        if (null == pid) {
            throw new IllegalArgumentException("pid can not be null");
        }

        Pageable pageable = PageRequest.of(page - 1, pageSize);
        List<RegionPo> poPage;
        if (all) {
            poPage = this.regionPoMapper.findByPid(pid, pageable);
        } else {
            poPage = this.regionPoMapper.findByPidEqualsAndStatusEquals(pid, Region.VALID, pageable);
        }
        return poPage.stream()
                .map(po -> this.build(po, Optional.ofNullable(null)))
                .collect(Collectors.toList());
    }
}
