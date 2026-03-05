package com.longhe.learn.oomall.payment.dao.bo;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.longhe.learn.javaee.core.copyfrom.CopyFrom;
import com.longhe.learn.javaee.core.copyfrom.CopyTo;
import com.longhe.learn.javaee.core.model.OOMallObject;
import com.longhe.learn.oomall.payment.dao.ChannelDao;
import com.longhe.learn.oomall.payment.mapper.po.AccountPo;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import lombok.ToString;
import lombok.extern.slf4j.Slf4j;

import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.*;

/**
 * 商铺收款账号
 */
@ToString(callSuper = true, doNotUseGetters = true)
@NoArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
@CopyFrom({AccountPo.class})
@CopyTo(AccountPo.class)
@Slf4j
public class Account extends OOMallObject implements Serializable {


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
     *  删除
     */
    @ToString.Exclude
    @JsonIgnore
    public static Byte DELETE = 2;

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
                    add(DELETE);
                }
            });
        }
    };

    /**
     * 状态和名称的对应
     */
    public static final Map<Byte, String> STATUSNAMES = new HashMap(){
        {
            put(VALID, "有效");
            put(INVALID, "无效");
            put(DELETE, "删除");
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

        if (!Objects.isNull(status) && !Objects.isNull(this.status)){
            Set<Byte> allowStatusSet = toStatus.get(this.status);
            if (!Objects.isNull(allowStatusSet)) {
                ret = allowStatusSet.contains(status);
            }
        }
        return ret;
    }


    /**
     * 商铺id
     */
    @Setter
    @Getter
    private Long shopId;

    /**
     * 子商户号
     */
    @Setter
    @Getter
    private String subMchid;

    /**
     * 状态
     */
    @Setter
    private Byte status;

    /**
     * 支付渠道
     */
    @ToString.Exclude
    @JsonIgnore
    private Channel channel;

    @Setter
    @Getter
    private Long channelId;

    @ToString.Exclude
    @JsonIgnore
    @Setter
    private ChannelDao channelDao;

    public Channel getChannel() {
        if (null == this.channel && null != this.channelDao) {
            log.debug("getChannel: channelId = {}", this.channelId);
            this.channel = this.channelDao.findById(this.channelId);
        }
        return this.channel;
    }


    @JsonIgnore
    public Byte getStatus() {
        log.debug("getStatus: this.status = {}", this.status);
        Channel channel = this.getChannel();
        if (Channel.VALID.equals(channel.getStatus())) {
            return this.status;
        }else{
            //支付渠道无效
            return Account.INVALID;
        }
    }

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
