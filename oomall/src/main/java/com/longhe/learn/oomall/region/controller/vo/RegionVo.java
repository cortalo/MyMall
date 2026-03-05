package com.longhe.learn.oomall.region.controller.vo;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.javaee.core.copyfrom.CopyFrom;
import com.longhe.learn.javaee.core.model.IdNameTypeVo;
import com.longhe.learn.javaee.core.validation.NewGroup;
import com.longhe.learn.oomall.region.dao.bo.Region;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * 地区的视图对象
 * 用于向前端返回值
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
@NoArgsConstructor
@CopyFrom({Region.class})
@Data
public class RegionVo {
    private Long id;
    private String name;
    private Byte status;
    private Byte level;
    private String shortName;
    private String mergerName;
    private String pinyin;
    private Double lng;
    private Double lat;
    private String areaCode;
    private String zipCode;
    private String cityCode;
    private IdNameTypeVo creator;
    private LocalDateTime gmtCreate;
    private LocalDateTime gmtModified;
    private IdNameTypeVo modifier;
}
