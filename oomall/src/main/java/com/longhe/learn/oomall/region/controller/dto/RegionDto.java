package com.longhe.learn.oomall.region.controller.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.javaee.core.copyfrom.CopyTo;
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
 * 地区数据传输对象，用于从前端接收数据
 */
@NoArgsConstructor
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
@CopyTo(Region.class)
@Data
public class RegionDto {
    @NotBlank(message = "地区名不能为空", groups = {NewGroup.class})
    private String name;
    @NotBlank(message = "地区简称不能为空", groups = {NewGroup.class})
    private String shortName;
    @NotBlank(message = "地区全称不能为空", groups = {NewGroup.class})
    private String mergerName;
    @NotBlank(message = "地区拼音不能为空", groups = {NewGroup.class})
    private String pinyin;
    @NotNull(message = "经度不能为空", groups = {NewGroup.class})
    private Double lng;
    @NotNull(message = "纬度不能为空", groups = {NewGroup.class})
    private Double lat;
    @NotBlank(message = "地区码不能为空", groups = {NewGroup.class})
    private String areaCode;
    @NotBlank(message = "邮政编码不能为空", groups = {NewGroup.class})
    private String zipCode;
    @NotBlank(message = "电话区号不能为空", groups = {NewGroup.class})
    private String cityCode;
}
