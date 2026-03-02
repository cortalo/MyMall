package com.longhe.learn.mymall.region.controller;

import com.longhe.learn.mymall.core.model.IdNameTypeVo;
import com.longhe.learn.mymall.core.model.ReturnObject;
import com.longhe.learn.mymall.region.dao.bo.Region;
import com.longhe.learn.mymall.region.service.RegionService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.stream.Collectors;

/**
 * 地区内部api控制器
 */
@RestController
@RequestMapping(value = "/internal", produces = "application/json;charset=UTF-8")
public class InternalRegionController {
    private static final Logger logger = LoggerFactory.getLogger(InternalRegionController.class);
    private final RegionService regionService;

    @Autowired
    public InternalRegionController(RegionService regionService) {
        this.regionService = regionService;
    }

    @GetMapping("/regions/{id}/parents")
    public ReturnObject getParentsRegions(@PathVariable Long id) {
        List<Region> ancestors = this.regionService.retrieveParentsRegionsById(id);
        return new ReturnObject(
                ancestors.stream().map(o -> IdNameTypeVo.builder().id(o.getId()).name(o.getName()).build())
                        .collect(Collectors.toList())
        );
    }
}
