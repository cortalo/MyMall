package com.longhe.learn.oomall.region.controller;

import com.longhe.learn.javaee.core.model.IdNameTypeVo;
import com.longhe.learn.javaee.core.model.ReturnObject;
import com.longhe.learn.oomall.region.controller.dto.RegionDto;
import com.longhe.learn.oomall.region.dao.bo.Region;
import com.longhe.learn.oomall.region.service.RegionService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import static com.longhe.learn.javaee.core.util.Common.cloneObj;

@RestController
@RequestMapping(produces = "application/json;charset=UTF-8")
public class UnAuthorizedController {

    private final Logger logger = LoggerFactory.getLogger(UnAuthorizedController.class);
    private final RegionService regionService;

    @Autowired
    public UnAuthorizedController(RegionService regionService) {
        this.regionService = regionService;
    }

    @GetMapping("/regions/{id}")
    public ReturnObject findRegionById(@PathVariable Long id) {
        Region region = this.regionService.findById(id);
        RegionDto dto = cloneObj(region, RegionDto.class);
        dto.setCreator(IdNameTypeVo.builder().id(region.getCreatorId()).name(region.getCreatorName()).build());
        dto.setModifier(IdNameTypeVo.builder().id(region.getModifierId()).name(region.getModifierName()).build());
        return new ReturnObject(dto);
    }

}
