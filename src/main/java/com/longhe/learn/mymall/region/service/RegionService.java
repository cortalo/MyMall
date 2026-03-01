package com.longhe.learn.mymall.region.service;

import com.longhe.learn.mymall.core.model.dto.UserDto;
import com.longhe.learn.mymall.region.dao.bo.Region;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Transactional
public class RegionService {

    /**
     * 创建新的子地区
     *
     * @param id region id
     * @param region 下级region对象
     * @param user 登录用户
     * @return 新region对象，带id
     */
    public Region createSubRegions(Long id, Region region, UserDto user) {
        return null;
    }

}
