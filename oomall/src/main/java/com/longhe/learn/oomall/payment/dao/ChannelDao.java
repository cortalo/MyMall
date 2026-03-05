package com.longhe.learn.oomall.payment.dao;

import com.longhe.learn.javaee.core.exception.BusinessException;
import com.longhe.learn.javaee.core.mapper.RedisUtil;
import com.longhe.learn.javaee.core.model.ReturnNo;
import com.longhe.learn.javaee.core.util.CloneFactory;
import com.longhe.learn.oomall.payment.dao.bo.Channel;
import com.longhe.learn.oomall.payment.mapper.ChannelPoMapper;
import com.longhe.learn.oomall.payment.mapper.po.ChannelPo;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Repository;

import java.util.Objects;
import java.util.Optional;

@Repository
@RequiredArgsConstructor
@Slf4j
public class ChannelDao {
    public static final String KEY = "C%d";

    private final RedisUtil redisUtil;
    private final ChannelPoMapper channelPoMapper;

    public Channel findById(Long id) {

        String key = String.format(KEY, id);
        Channel channel = (Channel) redisUtil.get(key);
        if (Objects.nonNull(channel)) {
            this.build(channel);
            return channel;
        } else {
            Optional<ChannelPo> optPo = this.channelPoMapper.findById(id);
            if (optPo.isEmpty()) {
                throw new BusinessException(ReturnNo.RESOURCE_ID_NOTEXIST, String.format(ReturnNo.RESOURCE_ID_NOTEXIST.getMessage(), "支付渠道", id));
            }
            return this.build(optPo.get(), Optional.of(key));
        }

    }

    private Channel build(ChannelPo po, Optional<String> redisKey) {
        Channel channel = CloneFactory.copy(new Channel(), po);
        redisKey.ifPresent(key -> redisUtil.set(key, channel, -1));
        this.build(channel);
        return channel;
    }

    private Channel build(Channel bo) {
        bo.setChannelDao(this);
        return bo;
    }
}
