package com.longhe.learn.oomall.payment.mapper;

import com.longhe.learn.oomall.payment.mapper.po.ChannelPo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ChannelPoMapper extends JpaRepository<ChannelPo, Long> {
}
