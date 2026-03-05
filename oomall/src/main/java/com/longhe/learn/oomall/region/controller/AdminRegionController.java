package com.longhe.learn.oomall.region.controller;

import com.longhe.learn.javaee.core.aop.Audit;
import com.longhe.learn.javaee.core.aop.LoginUser;
import com.longhe.learn.javaee.core.exception.BusinessException;
import com.longhe.learn.javaee.core.model.IdNameTypeVo;
import com.longhe.learn.javaee.core.model.ReturnNo;
import com.longhe.learn.javaee.core.model.ReturnObject;
import com.longhe.learn.javaee.core.model.UserToken;
import com.longhe.learn.javaee.core.util.CloneFactory;
import com.longhe.learn.javaee.core.validation.NewGroup;
import com.longhe.learn.oomall.region.controller.dto.RegionDto;
import com.longhe.learn.oomall.region.controller.vo.RegionVo;
import com.longhe.learn.oomall.region.dao.bo.Region;
import com.longhe.learn.oomall.region.service.RegionService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import static com.longhe.learn.javaee.core.model.Constants.PLATFORM;


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
    public ReturnObject createSubRegions(@PathVariable Long did, @PathVariable Long id, @LoginUser UserToken user,
                                         @Validated(NewGroup.class) @RequestBody RegionDto dto) {
        if (!PLATFORM.equals(did)) {
            throw new BusinessException(ReturnNo.RESOURCE_ID_OUTSCOPE, String.format(ReturnNo.RESOURCE_ID_OUTSCOPE.getMessage(), "地区", id, did));
        }
        Region region = CloneFactory.copy(new Region(), dto);
        Region newRegion = this.regionService.createSubRegions(id, region, user);
        IdNameTypeVo vo = IdNameTypeVo.builder().id(newRegion.getId()).name(newRegion.getName()).build();
        return new ReturnObject(ReturnNo.CREATED, vo);
    }

    @DeleteMapping("/regions/{id}")
    @Audit(departName = "shops")
    public ReturnObject deleteRegionById(@PathVariable Long did, @PathVariable Long id, @LoginUser UserToken user) {
        if (!PLATFORM.equals(did)) {
            throw new BusinessException(ReturnNo.RESOURCE_ID_OUTSCOPE, String.format(ReturnNo.RESOURCE_ID_OUTSCOPE.getMessage(), "地区", id, did));
        }
        this.regionService.deleteRegion(id, user);
        return new ReturnObject(ReturnNo.OK);
    }

}
