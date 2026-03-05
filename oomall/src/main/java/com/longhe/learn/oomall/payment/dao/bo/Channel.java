package com.longhe.learn.oomall.payment.dao.bo;


import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.javaee.core.copyfrom.CopyFrom;
import com.longhe.learn.javaee.core.copyfrom.CopyTo;
import com.longhe.learn.javaee.core.model.OOMallObject;
import com.longhe.learn.oomall.payment.dao.ChannelDao;
import com.longhe.learn.oomall.payment.mapper.po.ChannelPo;
import lombok.*;

import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

/**
 * 商户支付渠道
 */
@ToString(callSuper = true, doNotUseGetters = true)
@EqualsAndHashCode(callSuper = true)
@NoArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
@CopyFrom({ChannelPo.class})
@CopyTo(ChannelPo.class)
public class Channel extends OOMallObject implements Serializable {


    /**
     * 有效
     */
    @ToString.Exclude
    @JsonIgnore
    public static Byte VALID = 0;
    /**
     * 无效
     */
    @ToString.Exclude
    @JsonIgnore
    public static Byte INVALID = 1;

    /**
     * 状态和名称的对应
     */
    public static final Map<Byte, String> STATUSNAMES = new HashMap(){
        {
            put(VALID, "有效");
            put(INVALID, "无效");
        }
    };

    /**
     * 允许的状态迁移
     */
    private static final Map<Byte, Set<Byte>> toStatus = new HashMap<>(){
        {
            put(VALID, new HashSet<>(){
                {
                    add(INVALID);
                }
            });
            put(INVALID, new HashSet<>(){
                {
                    add(VALID);
                }
            });
        }
    };

    /**
     * 是否允许状态迁移
     * @author Ming Qiu
     * <p>
     * date: 2022-11-13 0:25
     * @param status
     * @return
     */
    public boolean allowStatus(Byte status){
        boolean ret = false;

        if (null != status && null != this.status){
            Set<Byte> allowStatusSet = toStatus.get(this.status);
            if (null != allowStatusSet) {
                ret = allowStatusSet.contains(status);
            }
        }
        return ret;
    }

    /**
     * 平台应用id
     */
    @Setter
    @Getter
    private String spAppid;

    /**
     * 渠道名称
     */
    @Setter
    @Getter
    private String name;

    /**
     * 平台商户号
     */
    @Setter
    @Getter
    private String spMchid;

    /**
     * 开始时间
     */
    @Setter
    @Getter
    private LocalDateTime beginTime;

    /**
     * 结束时间
     */
    @Setter
    @Getter
    private LocalDateTime endTime;

    /**
     *  状态
     */
    @Setter
    @Getter
    private Byte status;

    /**
     * 适配对象名
     */
    @Setter
    @Getter
    private String beanName;

    /**
     * 通知地址
     */
    @Setter
    @Getter
    private String notifyUrl;

    @ToString.Exclude
    @JsonIgnore
    @Setter
    private ChannelDao channelDao;


    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getCreatorId() {
        return creatorId;
    }

    public void setCreatorId(Long creatorId) {
        this.creatorId = creatorId;
    }

    public String getCreatorName() {
        return creatorName;
    }

    public void setCreatorName(String creatorName) {
        this.creatorName = creatorName;
    }

    public Long getModifierId() {
        return modifierId;
    }

    public void setModifierId(Long modifierId) {
        this.modifierId = modifierId;
    }

    public String getModifierName() {
        return modifierName;
    }

    public void setModifierName(String modifierName) {
        this.modifierName = modifierName;
    }

    public LocalDateTime getGmtCreate() {
        return gmtCreate;
    }

    public LocalDateTime getGmtModified() {
        return gmtModified;
    }

    @Override
    public void setGmtCreate(LocalDateTime gmtCreate) {
        this.gmtCreate = gmtCreate;
    }

    @Override
    public void setGmtModified(LocalDateTime gmtModified) {
        this.gmtModified = gmtModified;
    }
}
