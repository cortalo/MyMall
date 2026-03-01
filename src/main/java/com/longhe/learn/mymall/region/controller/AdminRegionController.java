package com.longhe.learn.mymall.region.controller;

import com.longhe.learn.mymall.core.aop.Audit;
import com.longhe.learn.mymall.core.aop.LoginUser;
import com.longhe.learn.mymall.core.exception.BusinessException;
import com.longhe.learn.mymall.core.model.ReturnNo;
import com.longhe.learn.mymall.core.model.ReturnObject;
import com.longhe.learn.mymall.core.model.dto.IdNameTypeDto;
import com.longhe.learn.mymall.core.model.dto.UserDto;
import com.longhe.learn.mymall.core.validation.NewGroup;
import com.longhe.learn.mymall.region.controller.vo.RegionVo;
import com.longhe.learn.mymall.region.dao.bo.Region;
import com.longhe.learn.mymall.region.service.RegionService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import static com.longhe.learn.mymall.core.model.Constants.PLATFORM;
import static com.longhe.learn.mymall.core.util.Common.cloneObj;


/**
 * 地区管理员控制器
 */
@RestController /*Restful的Controller对象*/
@RequestMapping(value = "/shops/{did}", produces = "application/json;charset=UTF-8")
public class AdminRegionController {
    private final Logger logger = LoggerFactory.getLogger(AdminRegionController.class);
    private final RegionService regionService;

    @Autowired
    public AdminRegionController(RegionService regionService) {
        this.regionService = regionService;
    }

    @PostMapping("/regions/{id}/subregions")
    @Audit(departName = "shops")
    public ReturnObject createSubRegions(@PathVariable Long did, @PathVariable Long id, @LoginUser UserDto user,
                                         @Validated(NewGroup.class) @RequestBody RegionVo vo) {
        if (!PLATFORM.equals(did)) {
            throw new BusinessException(ReturnNo.RESOURCE_ID_OUTSCOPE, String.format(ReturnNo.RESOURCE_ID_OUTSCOPE.getMessage(), "地区", id, did));
        }
        Region region = cloneObj(vo, Region.class);
        Region newRegion = this.regionService.createSubRegions(id, region, user);
        IdNameTypeDto dto = IdNameTypeDto.builder().id(newRegion.getId()).name(newRegion.getName()).build();
        return new ReturnObject(ReturnNo.CREATED, dto);
    }

}
