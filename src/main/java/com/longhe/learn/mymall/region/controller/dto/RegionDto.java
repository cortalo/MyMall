package com.longhe.learn.mymall.region.controller.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.mymall.core.model.dto.IdNameTypeDto;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@JsonInclude(JsonInclude.Include.NON_NULL)
@NoArgsConstructor
public class RegionDto {
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
    private IdNameTypeDto creator;
    private LocalDateTime gmtCreate;
    private LocalDateTime gmtModified;
    private IdNameTypeDto modifier;
}
