package com.longhe.learn.mymall.region.dao.bo;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.mymall.region.dao.RegionDao;
import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;
import lombok.Setter;
import lombok.ToString;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Serializable;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

@NoArgsConstructor
@AllArgsConstructor
@ToString(callSuper = true, doNotUseGetters = true)
@JsonInclude(JsonInclude.Include.NON_NULL)
public class Region implements Serializable {

    @ToString.Exclude
    @JsonIgnore
    private final static Logger logger = LoggerFactory.getLogger(Region.class);

    /**
     * 两种特殊id
     * 0 -- 最高级地区
     * -1 -- 不存在
     */
    @ToString.Exclude
    @JsonIgnore
    public static final Long TOP_ID = 0L;
    @ToString.Exclude
    @JsonIgnore
    public static final Long ROOT_PID = -1L;

    /**
     * 共三种状态
     */
    //有效
    @ToString.Exclude
    @JsonIgnore
    public static final Byte VALID = 0;
    //停用
    @ToString.Exclude
    @JsonIgnore
    public static final Byte SUSPENDED = 1;
    //废弃
    @ToString.Exclude
    @JsonIgnore
    public static final Byte ABANDONED = 2;

    @ToString.Exclude
    @JsonIgnore
    public static final Map<Byte, String> STATUSNAMES = new HashMap<>() {
        {
            put(VALID, "有效");
            put(SUSPENDED, "停用");
            put(ABANDONED, "废弃");
        }
    };

    /**
     * 允许的状态迁移
     */
    @JsonIgnore
    @ToString.Exclude
    private static final Map<Byte, Set<Byte>> toStatus = new HashMap<>() {
        {
            put(VALID, new HashSet<>() {
                {
                    add(SUSPENDED);
                    add(ABANDONED);
                }
            });
            put(SUSPENDED, new HashSet<>() {
                {
                    add(VALID);
                    add(ABANDONED);
                }
            });
        }
    };

    /**
     * 是否允许状态迁移
     */
    public boolean allowStatus(Byte status) {
        boolean ret = false;

        if (null != status && null != this.status) {
            Set<Byte> allowStatusSet = toStatus.get(this.status);
            if (null != allowStatusSet) {
                ret = allowStatusSet.contains(status);
            }
        }
        return ret;
    }

    /**
     * 获得当前状态名称
     */
    @JsonIgnore
    public String getStatusName() {
        return STATUSNAMES.get(this.status);
    }


    private Long pid;
    private Byte level;
    private String areaCode;
    private String zipCode;
    private String cityCode;
    private String name;
    private String shortName;
    private String mergerName;
    private String pinyin;
    private Double lng;
    private Double lat;
    private Byte status;

    @JsonIgnore
    @ToString.Exclude
    private Region parentRegion;

    @Setter
    @JsonIgnore
    @ToString.Exclude
    private RegionDao regionDao;

    public Long getPid() {
        return pid;
    }

    public void setPid(Long pid) {
        this.pid = pid;
    }

    public Byte getLevel() {
        return level;
    }

    public void setLevel(Byte level) {
        this.level = level;
    }

    public String getAreaCode() {
        return areaCode;
    }

    public void setAreaCode(String areaCode) {
        this.areaCode = areaCode;
    }

    public String getZipCode() {
        return zipCode;
    }

    public void setZipCode(String zipCode) {
        this.zipCode = zipCode;
    }

    public String getCityCode() {
        return cityCode;
    }

    public void setCityCode(String cityCode) {
        this.cityCode = cityCode;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getShortName() {
        return shortName;
    }

    public void setShortName(String shortName) {
        this.shortName = shortName;
    }

    public String getMergerName() {
        return mergerName;
    }

    public void setMergerName(String mergerName) {
        this.mergerName = mergerName;
    }

    public String getPinyin() {
        return pinyin;
    }

    public void setPinyin(String pinyin) {
        this.pinyin = pinyin;
    }

    public Double getLng() {
        return lng;
    }

    public void setLng(Double lng) {
        this.lng = lng;
    }

    public Double getLat() {
        return lat;
    }

    public void setLat(Double lat) {
        this.lat = lat;
    }

    public Byte getStatus() {
        return status;
    }

    public void setStatus(Byte status) {
        this.status = status;
    }
}
